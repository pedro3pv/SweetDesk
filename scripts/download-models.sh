#!/bin/bash
mkdir -p models/realcugan models/lsdir models/classifier

# RealCUGAN
if [ ! -f models/realcugan/realcugan-pro.onnx ]; then
  echo "Baixando RealCUGAN..."
  curl -L "https://huggingface.co/deepghs/imgutils-models/resolve/main/real_esrgan/RealESRGAN_x4plus_anime_6B.onnx" \
    -o models/realcugan/realcugan-pro.onnx
fi

# LSDIR
if [ ! -f models/lsdir/4xLSDIR.onnx ]; then
  echo "Baixando LSDIR..."
  curl -L "https://huggingface.co/wanesoft/faceswap_pack/resolve/main/lsdir_x4.onnx" \
    -o models/lsdir/4xLSDIR.onnx
fi

# Classifier
if [ ! -f models/classifier/model.onnx ]; then
  echo "Baixando Classifier..."
  curl -L "https://huggingface.co/deepghs/anime_real_cls/resolve/main/caformer_s36_v1.3_fixed/model.onnx" \
    -o models/classifier/model.onnx
fi
