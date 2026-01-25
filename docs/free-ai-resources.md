# FarmAgent: 100% Free AI Resources Guide
## Build Your AI Agent Without Spending a Dollar

This guide provides completely free datasets, models, APIs, and tools to build FarmAgent for the AIFEST 2026 hackathon.

---

## 1. Free Crop Disease Detection Datasets

### ✅ PlantVillage Dataset (FREE)
- **Link:** https://www.kaggle.com/datasets/emmarex/plantdisease
- **Size:** 54,309 images
- **Crops:** 14 crops (tomato, potato, corn, cassava, etc.)
- **Diseases:** 38 disease categories
- **License:** Creative Commons (Free for commercial use)
- **Download:** Free Kaggle account required

**How to use:**
```bash
# Install Kaggle CLI
pip install kaggle

# Download dataset
kaggle datasets download -d emmarex/plantdisease
unzip plantdisease.zip
```

### ✅ Cassava Leaf Disease Dataset (FREE)
- **Link:** https://www.kaggle.com/c/cassava-leaf-disease-classification
- **Size:** 21,367 images
- **Diseases:** 5 classes (CMD, CBSD, CGM, CBB, Healthy)
- **License:** Free for research and commercial
- **Perfect for:** Uganda-specific crops

### ✅ PlantDoc Dataset (FREE)
- **Link:** https://github.com/pratikkayal/PlantDoc-Dataset
- **Size:** 2,598 images
- **Species:** 13 plant species
- **Diseases:** 17 disease classes
- **License:** Creative Commons Attribution 4.0 International
- **Best for:** Real-world, non-lab conditions

### ✅ CCMT Dataset - Ghana Crops (FREE)
- **Link:** https://pmc.ncbi.nlm.nih.gov/articles/PMC10285554/
- **Size:** 24,881 raw images, 102,976 augmented
- **Crops:** Cassava, Maize, Cashew, Tomato
- **Perfect for:** African crop diseases
- **License:** Freely available for research

### ✅ Plant Leaf Disease Recognition Dataset (FREE)
- **Link:** https://data.mendeley.com/datasets/5g238dv4ht/1
- **Size:** 4,121 original, 20,956 augmented
- **Plants:** Gourd, Hibiscus, Papaya, Zucchini
- **Updated:** November 2024
- **License:** Free to use

---

## 2. Free AI/ML Models for Disease Detection

### ✅ Pre-trained Models from TensorFlow Hub (FREE)

**EfficientNet (Recommended)**
```python
import tensorflow_hub as hub

# Load pre-trained EfficientNet
model_url = "https://tfhub.dev/google/efficientnet/b3/feature-vector/1"
base_model = hub.KerasLayer(model_url, trainable=False)

# FREE - No API key needed
```

**MobileNetV2 (Lightweight)**
```python
from tensorflow.keras.applications import MobileNetV2

# Completely free
base_model = MobileNetV2(weights='imagenet', include_top=False)
```

**ResNet50**
```python
from tensorflow.keras.applications import ResNet50

base_model = ResNet50(weights='imagenet', include_top=False)
```

### ✅ PyTorch Pre-trained Models (FREE)
```python
import torchvision.models as models

# All free from torchvision
efficientnet = models.efficientnet_b3(pretrained=True)
resnet = models.resnet50(pretrained=True)
mobilenet = models.mobilenet_v2(pretrained=True)
```

---

## 3. Free Open-Source LLM Alternatives (No API Cost)

### Option 1: Ollama (Run Locally - 100% FREE)

**Installation:**
```bash
# macOS/Linux
curl -fsSL https://ollama.ai/install.sh | sh

# Windows: Download from ollama.ai

# Run any LLM locally for FREE
ollama pull llama3.2:3b       # 3B model (fast)
ollama pull mistral:7b        # 7B model (balanced)
ollama pull qwen2.5:7b        # 7B model (multilingual)
ollama pull gemma2:2b         # 2B model (ultra-fast)
```

**Use in Python:**
```python
import requests

def get_recommendation(disease, crop):
    response = requests.post('http://localhost:11434/api/generate', 
        json={
            "model": "llama3.2:3b",
            "prompt": f"Recommend treatment for {disease} in {crop} for Uganda farmers",
            "stream": False
        }
    )
    return response.json()['response']

# 100% FREE - runs on your computer
```

**Best Models for FarmAgent:**
- **llama3.2:3b** - Fast, good for mobile deployment
- **qwen2.5:7b** - Multilingual, supports local languages
- **mistral:7b** - Balanced performance
- **gemma2:2b** - Ultra-lightweight

### Option 2: HuggingFace Inference API (FREE Tier)

**Serverless Inference - FREE**
```python
from huggingface_hub import InferenceClient

client = InferenceClient()

# FREE - No API key for public models
response = client.text_generation(
    "Recommend organic pesticide for cassava brown streak disease",
    model="mistralai/Mistral-7B-Instruct-v0.2"
)

print(response)
```

**Free Limits:**
- Models < 10GB: Unlimited (with rate limits)
- Popular models: Supported even if > 10GB
- No credit card required

**Best FREE Models on HuggingFace:**
- `mistralai/Mistral-7B-Instruct-v0.2`
- `meta-llama/Llama-3.2-3B-Instruct`
- `google/gemma-2-2b-it`
- `Qwen/Qwen2.5-7B-Instruct`

### Option 3: Google Gemini FREE Tier

**Generous Free Tier**
```python
import google.generativeai as genai

# FREE API key from https://makersuite.google.com/app/apikey
genai.configure(api_key='YOUR_FREE_API_KEY')

model = genai.GenerativeModel('gemini-2.0-flash-exp')

response = model.generate_content(
    "Recommend treatment for maize disease in Uganda"
)

print(response.text)
```

**Free Limits:**
- 15 requests per minute
- 1,500 requests per day
- 1 million tokens per day
- **Perfect for hackathon!**

### Option 4: OpenRouter (FREE Tier)

**Access Multiple Models FREE**
```python
import requests

response = requests.post(
    "https://openrouter.ai/api/v1/chat/completions",
    headers={"Authorization": "Bearer YOUR_FREE_KEY"},
    json={
        "model": "google/gemini-2.0-flash-exp:free",  # FREE
        "messages": [
            {"role": "user", "content": "Treatment for cassava disease"}
        ]
    }
)
```

**Free Models on OpenRouter:**
- `google/gemini-2.0-flash-exp:free`
- `meta-llama/llama-3.2-3b-instruct:free`
- `qwen/qwen-2.5-7b-instruct:free`
- `mistralai/mistral-7b-instruct:free`

---

## 4. Free Weather APIs

### ✅ OpenWeather FREE Tier
- **URL:** https://openweathermap.org/api
- **Free Tier:** 1,000 calls/day, 60 calls/minute
- **Features:**
  - Current weather
  - 5-day forecast
  - Weather alerts
- **Perfect for hackathon!**

**Usage:**
```python
import requests

API_KEY = "your_free_key"  # Get from openweathermap.org
lat, lon = 0.3476, 32.5825  # Kampala

url = f"https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API_KEY}"
weather = requests.get(url).json()

# FREE 1000 calls/day
```

### ✅ Open-Meteo (100% FREE - No API Key)
- **URL:** https://open-meteo.com/
- **Completely FREE** - No registration needed
- **No rate limits** for non-commercial use
- **Features:**
  - Current weather
  - 7-day forecast
  - Historical data

**Usage:**
```python
import requests

url = "https://api.open-meteo.com/v1/forecast"
params = {
    "latitude": 0.3476,
    "longitude": 32.5825,
    "current_weather": True,
    "daily": ["temperature_2m_max", "precipitation_sum"]
}

weather = requests.get(url, params=params).json()
# 100% FREE - No API key needed!
```

---

## 5. Free Market Price Data

### ✅ WFP Price Database (FREE)
- **URL:** https://data.humdata.org/dataset/wfp-food-prices
- **Data:** Uganda food prices by market
- **Format:** CSV download
- **Updates:** Monthly
- **License:** Free and open

**Download:**
```python
import pandas as pd

# Direct download URL (FREE)
url = "https://data.humdata.org/dataset/wfp-food-prices-for-uganda"
prices = pd.read_csv(url)
```

### ✅ FAOSTAT (FREE)
- **URL:** http://www.fao.org/faostat
- **Data:** Agricultural statistics, production data
- **Free API:** Available
- **License:** Open data

---

## 6. Free SMS/Communication (Testing)

### ✅ Africa's Talking Sandbox (FREE)
- **URL:** https://africastalking.com/
- **Free Credits:** $1 on signup (100 SMS)
- **Sandbox:** Unlimited testing to test numbers

**Setup:**
```python
import africastalking

africastalking.initialize(
    username="sandbox",  # FREE sandbox
    api_key="your_sandbox_key"
)

sms = africastalking.SMS
sms.send("Test message", ["+254711XXXYYY"])  # FREE in sandbox
```

### ✅ Twilio Trial (FREE)
- **URL:** https://www.twilio.com/
- **Free Credits:** $15 trial credit
- **Good for:** Testing during hackathon

---

## 7. Free Development Tools

### ✅ Ollama (Local LLM Runtime - FREE)
- Run any open-source LLM locally
- No API costs ever
- Perfect for development

### ✅ LM Studio (FREE)
- **URL:** https://lmstudio.ai/
- Visual interface to run LLMs
- Download and run models locally
- Zero cost

### ✅ llama.cpp (FREE)
- **URL:** https://github.com/ggerganov/llama.cpp
- Run LLMs in pure C/C++
- Extremely efficient
- Free and open-source

### ✅ Jan.ai (FREE)
- **URL:** https://jan.ai/
- Desktop app for local LLMs
- 100% offline
- Privacy-focused

---

## 8. Free Image Storage

### ✅ Cloudinary FREE Tier
- **URL:** https://cloudinary.com/
- **Free:** 25 GB storage
- **Bandwidth:** 25 GB/month
- **Transformations:** 25,000/month

### ✅ ImgBB (FREE)
- **URL:** https://imgbb.com/
- Completely free image hosting
- API available
- No limits mentioned

### ✅ Supabase Storage (FREE)
- **URL:** https://supabase.com/
- **Free:** 1 GB storage
- **Bandwidth:** 2 GB
- **Perfect for** hackathon prototypes

---

## 9. Recommended FREE Stack for Hackathon

### Best Combination (100% FREE):

**For AI/NLP:**
1. **Ollama + Llama 3.2 3B** (Local, free forever)
   - OR **Google Gemini 2.0 Flash** (1,500 req/day free)
   - OR **HuggingFace Inference API** (Free tier)

**For Disease Detection:**
2. **PlantVillage Dataset** (54K images, free)
3. **EfficientNetB3** (Pre-trained, free)
4. **TensorFlow/PyTorch** (Free frameworks)

**For Weather:**
5. **Open-Meteo** (100% free, no API key)
   - OR **OpenWeather** (1,000/day free)

**For Market Prices:**
6. **WFP Data** (Free CSV downloads)
7. **Manual scraping** of Uganda Commodity Exchange

**For SMS:**
8. **Africa's Talking Sandbox** ($1 free credit)
9. **In-app notifications only** (Firebase free tier)

**For Storage:**
10. **Supabase** (1GB free)

---

## 10. Step-by-Step Setup (All FREE)

### Step 1: Install Ollama (FREE LLM)
```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Download a model (FREE)
ollama pull llama3.2:3b

# Start server
ollama serve
```

### Step 2: Download PlantVillage Dataset
```bash
# Install Kaggle CLI
pip install kaggle

# Set up Kaggle credentials (free account)
mkdir ~/.kaggle
# Download API key from kaggle.com/settings

# Download dataset (FREE)
kaggle datasets download -d emmarex/plantdisease
unzip plantdisease.zip
```

### Step 3: Train Disease Detection Model
```python
import tensorflow as tf
from tensorflow.keras.applications import EfficientNetB3

# Load pre-trained model (FREE)
base_model = EfficientNetB3(
    weights='imagenet',
    include_top=False,
    input_shape=(224, 224, 3)
)

# Add custom layers
model = tf.keras.Sequential([
    base_model,
    tf.keras.layers.GlobalAveragePooling2D(),
    tf.keras.layers.Dense(38, activation='softmax')  # 38 diseases
])

# Train on PlantVillage (FREE)
model.compile(optimizer='adam', loss='categorical_crossentropy')
# model.fit(train_data, ...)

# Save model (FREE)
model.save('disease_model.h5')
```

### Step 4: Set Up FREE LLM API
```python
import requests

# Option 1: Ollama (Local - FREE)
def get_treatment_local(disease):
    response = requests.post('http://localhost:11434/api/generate', 
        json={
            "model": "llama3.2:3b",
            "prompt": f"Recommend organic treatment for {disease} in Uganda",
            "stream": False
        }
    )
    return response.json()['response']

# Option 2: Gemini (1500 req/day FREE)
import google.generativeai as genai

genai.configure(api_key='YOUR_FREE_KEY')
model = genai.GenerativeModel('gemini-2.0-flash-exp')

def get_treatment_gemini(disease):
    response = model.generate_content(f"Treatment for {disease}")
    return response.text

# Option 3: HuggingFace (FREE)
from huggingface_hub import InferenceClient

client = InferenceClient()

def get_treatment_hf(disease):
    return client.text_generation(
        f"Recommend treatment: {disease}",
        model="mistralai/Mistral-7B-Instruct-v0.2"
    )
```

### Step 5: Get FREE Weather Data
```python
import requests

# Open-Meteo (100% FREE, no key needed)
def get_weather(lat, lon):
    url = "https://api.open-meteo.com/v1/forecast"
    params = {
        "latitude": lat,
        "longitude": lon,
        "current_weather": True
    }
    return requests.get(url, params=params).json()

weather = get_weather(0.3476, 32.5825)  # Kampala
# 100% FREE!
```

---

## 11. Cost Breakdown (All FREE for Hackathon)

| Resource | Free Tier | Cost for Hackathon |
|----------|-----------|-------------------|
| **Datasets** | PlantVillage, Cassava, PlantDoc | $0 |
| **AI Models** | Ollama + Llama 3.2 3B | $0 |
| **OR** Gemini 2.0 Flash | 1,500 req/day | $0 |
| **OR** HuggingFace API | Unlimited (rate limited) | $0 |
| **Weather API** | Open-Meteo unlimited | $0 |
| **Market Data** | WFP CSV downloads | $0 |
| **SMS Testing** | Africa's Talking $1 credit | $0 |
| **Image Storage** | Supabase 1GB | $0 |
| **Development** | Ollama, VS Code, etc | $0 |
| **Total** | | **$0** |

---

## 12. Pro Tips for Free Tier Usage

### Maximize Gemini Free Tier
```python
# Batch requests to save quota
def batch_recommendations(diseases):
    prompt = "Provide treatments for: " + ", ".join(diseases)
    response = model.generate_content(prompt)
    return response.text

# Cache responses
import functools

@functools.lru_cache(maxsize=1000)
def get_cached_recommendation(disease):
    return model.generate_content(f"Treatment: {disease}").text
```

### Use Ollama for Development
```bash
# During development: Use Ollama (local, unlimited)
ollama pull llama3.2:3b

# For demo: Switch to Gemini (faster, cloud-based)
# Both FREE!
```

### Optimize Image Storage
```python
from PIL import Image
import io

def compress_image(image_path, max_size_mb=1):
    """Compress images before upload to save storage"""
    img = Image.open(image_path)
    img.thumbnail((800, 800))  # Resize
    
    buffer = io.BytesIO()
    img.save(buffer, format='JPEG', quality=85, optimize=True)
    
    return buffer.getvalue()
```

---

## 13. Fallback Options (Still Free)

### If Gemini Quota Runs Out:
1. Switch to **Ollama** (local, unlimited)
2. Use **HuggingFace Inference API**
3. Try **OpenRouter free tier**
4. Use **Groq free tier** (ultra-fast)

### If Weather API Limits Hit:
1. Use **Open-Meteo** (no limits)
2. Cache weather data (updates every 3 hours)
3. Use static weather data for demo

---

## 14. Quick Start Script (All FREE)

```bash
#!/bin/bash

echo "Setting up FarmAgent with 100% FREE resources..."

# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Download LLM
ollama pull llama3.2:3b

# Install Python dependencies
pip install tensorflow pillow opencv-python pandas requests google-generativeai huggingface-hub

# Download PlantVillage dataset
kaggle datasets download -d emmarex/plantdisease
unzip plantdisease.zip -d datasets/

echo "Setup complete! All resources are FREE."
echo "Start Ollama: ollama serve"
echo "Train model: python train_model.py"
```

---

## 15. Summary: Your FREE AI Stack

✅ **Disease Detection:**
- PlantVillage Dataset (54K images)
- EfficientNetB3 (pre-trained)
- TensorFlow (framework)

✅ **Recommendations:**
- Ollama + Llama 3.2 3B (local, unlimited)
- OR Gemini 2.0 Flash (1,500/day)
- OR HuggingFace Inference (free tier)

✅ **Weather:**
- Open-Meteo (unlimited, no key)
- OR OpenWeather (1,000/day)

✅ **Market Prices:**
- WFP Data (free downloads)

✅ **Testing:**
- Africa's Talking ($1 free credit)
- Local development (unlimited)

**Total Cost: $0**
**Perfect for:** AIFEST 2026 Hackathon

---

## Resources Links Summary

### Datasets (All FREE):
- https://www.kaggle.com/datasets/emmarex/plantdisease
- https://www.kaggle.com/c/cassava-leaf-disease-classification
- https://github.com/pratikkayal/PlantDoc-Dataset
- https://data.mendeley.com/datasets/5g238dv4ht/1

### LLMs (All FREE):
- https://ollama.ai/ (Local)
- https://makersuite.google.com/ (Gemini)
- https://huggingface.co/inference-api
- https://openrouter.ai/ (Free tier)

### APIs (All FREE):
- https://open-meteo.com/ (Weather)
- https://openweathermap.org/ (Weather)
- https://data.humdata.org/ (Market prices)
- https://africastalking.com/ (SMS - $1 credit)

### Tools (All FREE):
- https://ollama.ai/ (LLM runtime)
- https://lmstudio.ai/ (LLM UI)
- https://jan.ai/ (Offline LLM)
- https://supabase.com/ (Storage)

---

**You can build the entire FarmAgent AI without spending a single dollar!** 🎉