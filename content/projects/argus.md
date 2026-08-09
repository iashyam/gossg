---
title: "Argus (Visual Intelligence)"
date: "2024-05-15"
tags: ["vision-models", "pytorch", "flask", "image-classification"]
image: "assets/argus.png"
link: "https://github.com/iashyam/whatisthis"
description: "A from-scratch implementation of vision models for image classification (SimpleCNN, MobileNet V2) in PyTorch, served via a Flask API and Docker."
---

A from-scratch implementation of vision models for image classification, served via a Flask API.

### Key Features
- **From-Scratch Architectures**: Custom CNN models and MobileNet V2 implementation built completely from scratch in PyTorch.
- **ImageNet Classification**: Classified into 1000 categories without pre-trained model dependencies at inference time.
- **Robust Pipeline**: Includes custom ETL pipelines (`ImageDataset`), a training loop module (`Trainer` class), MLflow experiment tracking, and unit tests via GitHub Actions.
- **Production Ready**: Containerized with Docker and served via Cloudflare Tunnels.
