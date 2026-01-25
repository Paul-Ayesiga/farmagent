#!/usr/bin/env python3
"""
Train Disease Classification Model

Fine-tunes EfficientNetB3 on PlantVillage dataset for crop disease detection.
"""

import os
import argparse
from pathlib import Path


def download_dataset(dataset_dir: str = "./datasets/plantvillage"):
    """Download PlantVillage dataset from Kaggle."""
    import subprocess
    
    dataset_path = Path(dataset_dir)
    
    if dataset_path.exists() and any(dataset_path.iterdir()):
        print(f"✅ Dataset already exists at {dataset_dir}")
        return dataset_path
    
    print("📥 Downloading PlantVillage dataset from Kaggle...")
    print("   (Requires Kaggle API credentials in ~/.kaggle/kaggle.json)")
    
    try:
        # Create directory
        dataset_path.mkdir(parents=True, exist_ok=True)
        
        # Download using Kaggle CLI
        subprocess.run([
            "kaggle", "datasets", "download",
            "-d", "emmarex/plantdisease",
            "-p", str(dataset_path),
            "--unzip"
        ], check=True)
        
        print(f"✅ Dataset downloaded to {dataset_dir}")
        return dataset_path
        
    except subprocess.CalledProcessError as e:
        print(f"❌ Failed to download dataset: {e}")
        print("   Please download manually from:")
        print("   https://www.kaggle.com/datasets/emmarex/plantdisease")
        raise


def train_model(
    dataset_dir: str = "./datasets/plantvillage",
    output_dir: str = "./models",
    epochs: int = 10,
    batch_size: int = 32,
    learning_rate: float = 0.001,
    input_size: int = 224,
):
    """Train the disease classification model."""
    import tensorflow as tf
    from tensorflow.keras.applications import EfficientNetB3
    from tensorflow.keras.layers import Dense, GlobalAveragePooling2D, Dropout
    from tensorflow.keras.models import Sequential
    from tensorflow.keras.optimizers import Adam
    from tensorflow.keras.callbacks import ModelCheckpoint, EarlyStopping, ReduceLROnPlateau
    
    print("🧠 Training Disease Classification Model")
    print(f"   Dataset: {dataset_dir}")
    print(f"   Output: {output_dir}")
    print(f"   Epochs: {epochs}")
    print(f"   Batch size: {batch_size}")
    
    # Find the actual data directory
    data_dir = Path(dataset_dir)
    
    # PlantVillage structure: plantvillage/color/
    for subdir in ["PlantVillage", "plantvillage", "color", "Plant_leave_diseases_dataset_with_augmentation"]:
        possible_path = data_dir / subdir
        if possible_path.exists():
            data_dir = possible_path
    
    print(f"   Data directory: {data_dir}")
    
    # Load dataset
    print("📂 Loading dataset...")
    
    train_ds = tf.keras.preprocessing.image_dataset_from_directory(
        data_dir,
        validation_split=0.2,
        subset="training",
        seed=42,
        image_size=(input_size, input_size),
        batch_size=batch_size,
        label_mode="categorical",
    )
    
    val_ds = tf.keras.preprocessing.image_dataset_from_directory(
        data_dir,
        validation_split=0.2,
        subset="validation",
        seed=42,
        image_size=(input_size, input_size),
        batch_size=batch_size,
        label_mode="categorical",
    )
    
    # Get class names
    class_names = train_ds.class_names
    num_classes = len(class_names)
    print(f"   Found {num_classes} disease classes")
    
    # Save class names
    output_path = Path(output_dir)
    output_path.mkdir(parents=True, exist_ok=True)
    
    with open(output_path / "class_names.txt", "w") as f:
        for name in class_names:
            f.write(f"{name}\n")
    
    # Optimize dataset
    AUTOTUNE = tf.data.AUTOTUNE
    train_ds = train_ds.cache().shuffle(1000).prefetch(buffer_size=AUTOTUNE)
    val_ds = val_ds.cache().prefetch(buffer_size=AUTOTUNE)
    
    # Data augmentation
    data_augmentation = tf.keras.Sequential([
        tf.keras.layers.RandomFlip("horizontal"),
        tf.keras.layers.RandomRotation(0.2),
        tf.keras.layers.RandomZoom(0.2),
        tf.keras.layers.RandomContrast(0.2),
    ])
    
    # Build model
    print("🏗️ Building model...")
    
    # Load pre-trained EfficientNetB3
    base_model = EfficientNetB3(
        weights="imagenet",
        include_top=False,
        input_shape=(input_size, input_size, 3),
    )
    
    # Freeze base model
    base_model.trainable = False
    
    # Build full model
    model = Sequential([
        # Preprocessing
        tf.keras.layers.Rescaling(1./255),
        data_augmentation,
        
        # Base model
        base_model,
        
        # Classification head
        GlobalAveragePooling2D(),
        Dropout(0.3),
        Dense(512, activation="relu"),
        Dropout(0.3),
        Dense(num_classes, activation="softmax"),
    ])
    
    # Compile
    model.compile(
        optimizer=Adam(learning_rate=learning_rate),
        loss="categorical_crossentropy",
        metrics=["accuracy"],
    )
    
    model.summary()
    
    # Callbacks
    callbacks = [
        ModelCheckpoint(
            str(output_path / "disease_classifier.h5"),
            monitor="val_accuracy",
            save_best_only=True,
            verbose=1,
        ),
        EarlyStopping(
            monitor="val_loss",
            patience=5,
            restore_best_weights=True,
            verbose=1,
        ),
        ReduceLROnPlateau(
            monitor="val_loss",
            factor=0.2,
            patience=3,
            min_lr=1e-7,
            verbose=1,
        ),
    ]
    
    # Train
    print("🚀 Starting training...")
    
    history = model.fit(
        train_ds,
        validation_data=val_ds,
        epochs=epochs,
        callbacks=callbacks,
    )
    
    # Fine-tune (unfreeze some layers)
    print("🔧 Fine-tuning...")
    
    base_model.trainable = True
    
    # Freeze early layers, train later layers
    fine_tune_at = len(base_model.layers) - 30
    for layer in base_model.layers[:fine_tune_at]:
        layer.trainable = False
    
    # Recompile with lower learning rate
    model.compile(
        optimizer=Adam(learning_rate=learning_rate / 10),
        loss="categorical_crossentropy",
        metrics=["accuracy"],
    )
    
    # Continue training
    history_fine = model.fit(
        train_ds,
        validation_data=val_ds,
        epochs=epochs // 2,
        initial_epoch=len(history.epoch),
        callbacks=callbacks,
    )
    
    # Save final model
    model.save(str(output_path / "disease_classifier_final.h5"))
    
    # Evaluate
    print("📊 Evaluating model...")
    loss, accuracy = model.evaluate(val_ds)
    print(f"   Validation Loss: {loss:.4f}")
    print(f"   Validation Accuracy: {accuracy:.4f}")
    
    print(f"✅ Model saved to {output_path / 'disease_classifier.h5'}")
    
    return model, history


def main():
    parser = argparse.ArgumentParser(description="Train disease classification model")
    parser.add_argument(
        "--dataset", 
        type=str, 
        default="./datasets/plantvillage",
        help="Path to PlantVillage dataset"
    )
    parser.add_argument(
        "--output",
        type=str,
        default="./models",
        help="Output directory for trained model"
    )
    parser.add_argument(
        "--epochs",
        type=int,
        default=10,
        help="Number of training epochs"
    )
    parser.add_argument(
        "--batch-size",
        type=int,
        default=32,
        help="Training batch size"
    )
    parser.add_argument(
        "--learning-rate",
        type=float,
        default=0.001,
        help="Initial learning rate"
    )
    parser.add_argument(
        "--download",
        action="store_true",
        help="Download dataset from Kaggle"
    )
    
    args = parser.parse_args()
    
    if args.download:
        download_dataset(args.dataset)
    
    train_model(
        dataset_dir=args.dataset,
        output_dir=args.output,
        epochs=args.epochs,
        batch_size=args.batch_size,
        learning_rate=args.learning_rate,
    )


if __name__ == "__main__":
    main()
