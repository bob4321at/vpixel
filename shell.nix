{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    opencv
    pkg-config
    glib
    glibc.dev
    stdenv.cc
    # Keep X11/media deps if needed
    xorg.libXrandr
    xorg.libXinerama
    xorg.libXcursor
    xorg.libXi
    xorg.libXxf86vm
    xorg.libX11
    xorg.libXext

    gtk3
    ffmpeg

    gcc

    python311
    python311Packages.pip
    python311Packages.numpy
    python311Packages.opencv4

    libGL

    mesa
    libglvnd
    xorg.libxcb
    alsa-lib
    dbus
    bear
    clang-tools
  ];

  LD_LIBRARY_PATH = with pkgs; lib.makeLibraryPath [
    mesa
    libglvnd
    xorg.libX11
    xorg.libXext
    alsa-lib
  ];

  OPENCV4_CFLAGS = "-I${pkgs.opencv}/include/opencv4";
  OPENCV4_LIBS = "-L${pkgs.opencv}/lib -lopencv_core -lopencv_imgproc -lopencv_videoio -lopencv_highgui";

  shellHook = ''
    export CGO_ENABLED=1
    export CC=${pkgs.stdenv.cc}/bin/cc
    export CXX=${pkgs.stdenv.cc}/bin/c++
  '';
}
