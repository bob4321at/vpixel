# tracker.py
import cv2
import mediapipe as mp
import json
import sys
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from socketserver import ThreadingMixIn
import time

# Global variable to hold latest tracking data
latest_data = {
    "left": None,
    "right": None,
    "mouth": {
        "left_corner": None,
        "right_corner": None,
        "upper_lip": None,
        "lower_lip": None
    },
    "timestamp": time.time()
}
data_lock = threading.Lock()

# MediaPipe setup
mp_face_mesh = mp.solutions.face_mesh
face_mesh = mp_face_mesh.FaceMesh(
    static_image_mode=False,
    max_num_faces=1,
    refine_landmarks=True,
    min_detection_confidence=0.5
)

LEFT_IRIS = [474, 475, 476, 477]
RIGHT_IRIS = [469, 470, 471, 472]
MOUTH_POINTS = {
    "left_corner": 78,
    "right_corner": 308,
    "upper_lip": 13,
    "lower_lip": 14
}

def tracking_loop():
    global latest_data
    cap = cv2.VideoCapture(0)
    if not cap.isOpened():
        error_msg = {"error": "Cannot open camera", "timestamp": time.time()}
        with data_lock:
            latest_data = error_msg
        print("❌ Camera failed to open", file=sys.stderr)
        return

    cap.set(cv2.CAP_PROP_FRAME_WIDTH, 320)
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 240)

    try:
        while True:
            ret, frame = cap.read()
            if not ret:
                break

            rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
            results = face_mesh.process(rgb)

            output = {
                "left": None,
                "right": None,
                "mouth": {
                    "left_corner": None,
                    "right_corner": None,
                    "upper_lip": None,
                    "lower_lip": None
                },
                "timestamp": time.time()
            }

            if results.multi_face_landmarks:
                landmarks = results.multi_face_landmarks[0].landmark
                h, w = frame.shape[0], frame.shape[1]

                left_x = sum(landmarks[i].x for i in LEFT_IRIS) / 4
                left_y = sum(landmarks[i].y for i in LEFT_IRIS) / 4
                output["left"] = [left_x * w, left_y * h]

                right_x = sum(landmarks[i].x for i in RIGHT_IRIS) / 4
                right_y = sum(landmarks[i].y for i in RIGHT_IRIS) / 4
                output["right"] = [right_x * w, right_y * h]

                for name, idx in MOUTH_POINTS.items():
                    lm = landmarks[idx]
                    output["mouth"][name] = [lm.x * w, lm.y * h]

            with data_lock:
                latest_data = output

    except Exception as e:
        error_msg = {"error": str(e), "timestamp": time.time()}
        with data_lock:
            latest_data = error_msg
        print(f"❌ Tracking error: {e}", file=sys.stderr)
    finally:
        cap.release()

class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
    """Handle requests in separate threads."""
    pass

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/' or self.path == '/data':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.send_header('Access-Control-Allow-Origin', '*')  # Allow CORS
            self.end_headers()

            with data_lock:
                response = latest_data

            self.wfile.write(json.dumps(response).encode('utf-8'))
        else:
            self.send_response(404)
            self.end_headers()

    def log_message(self, format, *args):
        # Optional: suppress console logs
        return

if __name__ == "__main__":
    # Start tracking in background thread
    tracker_thread = threading.Thread(target=tracking_loop, daemon=True)
    tracker_thread.start()

    # Give camera a moment to initialize
    time.sleep(1.0)

    server = ThreadedHTTPServer(('localhost', 8080), Handler)
    print("🚀 Face tracking server running on http://localhost:8080/")
    print("Press Ctrl+C to stop.")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n🛑 Shutting down...")
        server.shutdown()
