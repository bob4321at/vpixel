import cv2
import mediapipe as mp
from mediapipe.tasks import python
from mediapipe.tasks.python import vision
import json
import sys
import threading
import time
from http.server import BaseHTTPRequestHandler, HTTPServer
from socketserver import ThreadingMixIn

# =============================================================================
# CONFIGURATION
# =============================================================================
MODEL_PATH = "face_landmarker.task"  # Download from: https://goo.gle/3ey8zqx

# Landmark indices (478-landmark model)
LEFT_EYE_TOP = 159
LEFT_EYE_BOTTOM = 145
RIGHT_EYE_TOP = 386
RIGHT_EYE_BOTTOM = 374
MOUTH_POINTS = {
    "left_corner": 78,
    "right_corner": 308,
    "upper_lip": 13,
    "lower_lip": 14
}

# Global state
latest_data = {
    "left_eye": {"top": None, "bottom": None},
    "right_eye": {"top": None, "bottom": None},
    "mouth": {"left_corner": None, "right_corner": None, "upper_lip": None, "lower_lip": None},
    "timestamp": time.time(),
    "face_detected": False
}
data_lock = threading.Lock()
stop_event = threading.Event()


def tracking_loop():
    """Synchronous camera capture loop using VIDEO mode."""
    global latest_data
    
    # Configure for CPU-only (avoids EGL/GPU warnings on headless systems)
    base_options = python.BaseOptions(
        model_asset_path=MODEL_PATH,
        delegate=python.BaseOptions.Delegate.CPU  # Force CPU to avoid GPU init issues
    )
    options = vision.FaceLandmarkerOptions(
        base_options=base_options,
        running_mode=vision.RunningMode.VIDEO,  # Synchronous mode
        num_faces=1,
        min_face_detection_confidence=0.5,
        min_face_presence_confidence=0.5,
        min_tracking_confidence=0.5
    )
    
    try:
        landmarker = vision.FaceLandmarker.create_from_options(options)
        print("✅ FaceLandmarker loaded successfully")
    except RuntimeError as e:
        error_msg = {
            "error": f"Model load failed: {str(e)}",
            "hint": "Download face_landmarker.task from https://goo.gle/3ey8zqx",
            "timestamp": time.time()
        }
        with data_lock:
            latest_data = error_msg
        print(f"❌ {error_msg['error']}", file=sys.stderr)
        return

    cap = cv2.VideoCapture(0)
    if not cap.isOpened():
        error_msg = {"error": "Cannot open camera", "timestamp": time.time()}
        with data_lock:
            latest_data = error_msg
        print("❌ Camera failed to open", file=sys.stderr)
        landmarker.close()
        return

    # Set lower resolution for better performance
    cap.set(cv2.CAP_PROP_FRAME_WIDTH, 320)
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 240)
    
    frame_width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    frame_height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
    print(f"📷 Camera initialized: {frame_width}x{frame_height}")

    try:
        while not stop_event.is_set():
            ret, frame = cap.read()
            if not ret:
                print("⚠️ Frame capture failed")
                continue

            # Convert BGR → RGB for MediaPipe
            rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
            mp_image = mp.Image(image_format=mp.ImageFormat.SRGB, data=rgb_frame)
            timestamp_ms = int(time.time() * 1000)

            # Synchronous detection (VIDEO mode)
            result = landmarker.detect_for_video(mp_image, timestamp_ms)
            
            # Build output structure
            output = {
                "left_eye": {"top": None, "bottom": None},
                "right_eye": {"top": None, "bottom": None},
                "mouth": {"left_corner": None, "right_corner": None, "upper_lip": None, "lower_lip": None},
                "timestamp": time.time(),
                "face_detected": False
            }

            if result.face_landmarks and len(result.face_landmarks) > 0:
                landmarks = result.face_landmarks[0]  # First face
                output["face_detected"] = True
                
                # Helper to safely extract pixel coordinates
                def get_pixel(idx):
                    if idx < len(landmarks):
                        lm = landmarks[idx]
                        return [lm.x * frame_width, lm.y * frame_height]
                    return None
                
                # Extract eyelid points
                output["left_eye"]["top"] = get_pixel(LEFT_EYE_TOP)
                output["left_eye"]["bottom"] = get_pixel(LEFT_EYE_BOTTOM)
                output["right_eye"]["top"] = get_pixel(RIGHT_EYE_TOP)
                output["right_eye"]["bottom"] = get_pixel(RIGHT_EYE_BOTTOM)
                
                # Extract mouth points
                for name, idx in MOUTH_POINTS.items():
                    output["mouth"][name] = get_pixel(idx)

            # Thread-safe update
            with data_lock:
                latest_data = output

            # Optional: small delay to reduce CPU usage
            time.sleep(0.01)

    except Exception as e:
        error_msg = {"error": f"Tracking error: {str(e)}", "timestamp": time.time()}
        with data_lock:
            latest_data = error_msg
        print(f"❌ {error_msg['error']}", file=sys.stderr)
    finally:
        cap.release()
        landmarker.close()
        print("🔚 Resources cleaned up")


class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
    daemon_threads = True
    allow_reuse_address = True


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path in ['/', '/data']:
            with data_lock:
                response = latest_data.copy()
            
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.send_header('Access-Control-Allow-Origin', '*')
            self.end_headers()
            self.wfile.write(json.dumps(response).encode('utf-8'))
        else:
            self.send_response(404)
            self.end_headers()
    
    def log_message(self, format, *args):
        pass  # Suppress logs


if __name__ == "__main__":
    print("🚀 MediaPipe Face Landmarker (Tasks API - VIDEO mode)")
    print(f"📦 Model: {MODEL_PATH}")
    print("🔗 API: http://localhost:8080/data")
    print("💡 Tip: Download model if missing:")
    print("   curl -O https://storage.googleapis.com/mediapipe-models/face_landmarker/face_landmarker/float16/latest/face_landmarker.task\n")
    
    # Start tracking in background thread
    tracker = threading.Thread(target=tracking_loop, daemon=True)
    tracker.start()
    
    time.sleep(1.5)  # Wait for camera/model init
    
    # Start server
    server = ThreadedHTTPServer(('localhost', 8080), Handler)
    print("✅ Server ready. Press Ctrl+C to exit.\n")
    
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n🛑 Shutting down...")
        stop_event.set()
        tracker.join(timeout=2)
        server.shutdown()
