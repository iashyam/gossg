---
title: "Athena"
date: "2024-03-24"
tags: ["sentiment-analysis", "pytorch", "deep-learning", "huggingface"]
image: "assets/athena.png"
link: "https://github.com/iashyam/athena"
description: "An end-to-end PyTorch training pipeline for multiclass sentiment analysis featuring text embeddings caching, MLflow tracking, and ONNX export."
---

Athena is an end-to-end PyTorch training pipeline for multiclass sentiment analysis. It features automatic dataset downloading, text embeddings caching using HuggingFace's `SentenceTransformers`, model training, evaluation, and comprehensive experiment tracking via MLflow.

### Key Features
- **Automated Data Pipeline**: Automatically downloads and prepares raw dataset CSV files from HuggingFace.
- **Fast Text Embeddings**: Leverages `sentence-transformers` models (like `all-MiniLM-L6-v2`) to compute representations and caches them locally as PyTorch tensors (`.pt`) for incredibly fast iteration.
- **Config-Driven Training**: All major hyperparameters, data paths, and model configurations are extracted into a central `parameter.yaml` file.
- **MLflow & Databricks Integration**: Automatically tracks experiments, logs metrics, and registers models.
- **ONNX Export**: Trained PyTorch models are seamlessly converted and exported to ONNX format.
