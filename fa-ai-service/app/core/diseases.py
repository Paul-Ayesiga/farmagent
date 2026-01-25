"""
Disease Classes for FarmAgent

Based on PlantVillage dataset + African crop diseases
"""

# PlantVillage disease classes (38 classes)
PLANTVILLAGE_CLASSES = [
    "Apple___Apple_scab",
    "Apple___Black_rot",
    "Apple___Cedar_apple_rust",
    "Apple___healthy",
    "Blueberry___healthy",
    "Cherry_(including_sour)___Powdery_mildew",
    "Cherry_(including_sour)___healthy",
    "Corn_(maize)___Cercospora_leaf_spot Gray_leaf_spot",
    "Corn_(maize)___Common_rust_",
    "Corn_(maize)___Northern_Leaf_Blight",
    "Corn_(maize)___healthy",
    "Grape___Black_rot",
    "Grape___Esca_(Black_Measles)",
    "Grape___Leaf_blight_(Isariopsis_Leaf_Spot)",
    "Grape___healthy",
    "Orange___Haunglongbing_(Citrus_greening)",
    "Peach___Bacterial_spot",
    "Peach___healthy",
    "Pepper,_bell___Bacterial_spot",
    "Pepper,_bell___healthy",
    "Potato___Early_blight",
    "Potato___Late_blight",
    "Potato___healthy",
    "Raspberry___healthy",
    "Soybean___healthy",
    "Squash___Powdery_mildew",
    "Strawberry___Leaf_scorch",
    "Strawberry___healthy",
    "Tomato___Bacterial_spot",
    "Tomato___Early_blight",
    "Tomato___Late_blight",
    "Tomato___Leaf_Mold",
    "Tomato___Septoria_leaf_spot",
    "Tomato___Spider_mites Two-spotted_spider_mite",
    "Tomato___Target_Spot",
    "Tomato___Tomato_Yellow_Leaf_Curl_Virus",
    "Tomato___Tomato_mosaic_virus",
    "Tomato___healthy",
]

# Cassava diseases (from Kaggle cassava dataset)
CASSAVA_CLASSES = [
    "Cassava___Cassava_Mosaic_Disease",
    "Cassava___Cassava_Brown_Streak_Disease",
    "Cassava___Cassava_Green_Mottle",
    "Cassava___Cassava_Bacterial_Blight",
    "Cassava___healthy",
]

# Combined Uganda-relevant crop diseases
UGANDA_CROP_CLASSES = {
    "maize": [
        "Maize___Gray_Leaf_Spot",
        "Maize___Common_Rust",
        "Maize___Northern_Leaf_Blight",
        "Maize___Maize_Streak_Virus",
        "Maize___Fall_Armyworm",
        "Maize___healthy",
    ],
    "cassava": [
        "Cassava___Mosaic_Disease",
        "Cassava___Brown_Streak_Disease",
        "Cassava___Green_Mottle",
        "Cassava___Bacterial_Blight",
        "Cassava___healthy",
    ],
    "beans": [
        "Beans___Angular_Leaf_Spot",
        "Beans___Bean_Rust",
        "Beans___Common_Blight",
        "Beans___Anthracnose",
        "Beans___healthy",
    ],
    "tomato": [
        "Tomato___Bacterial_Spot",
        "Tomato___Early_Blight",
        "Tomato___Late_Blight",
        "Tomato___Leaf_Mold",
        "Tomato___Septoria_Leaf_Spot",
        "Tomato___Yellow_Leaf_Curl_Virus",
        "Tomato___Mosaic_Virus",
        "Tomato___healthy",
    ],
    "banana": [
        "Banana___Black_Sigatoka",
        "Banana___Panama_Disease",
        "Banana___Banana_Bunchy_Top_Virus",
        "Banana___Banana_Bacterial_Wilt",
        "Banana___healthy",
    ],
    "passion_fruit": [
        "PassionFruit___Woodiness_Virus",
        "PassionFruit___Brown_Spot",
        "PassionFruit___Fusarium_Wilt",
        "PassionFruit___healthy",
    ],
}

# Disease info with treatments
DISEASE_INFO = {
    "Maize___Fall_Armyworm": {
        "crop": "Maize",
        "disease": "Fall Armyworm",
        "severity_levels": ["low", "moderate", "severe"],
        "treatments": {
            "organic": [
                "Neem oil spray (20ml per liter water)",
                "Bacillus thuringiensis (Bt) spray",
                "Handpicking larvae in early stages",
                "Wood ash application",
            ],
            "chemical": [
                "Ampligo 150ZC (20ml per 20L water)",
                "Match 050EC (30ml per 20L water)",
                "Radiant SC (10ml per 20L water)",
            ],
        },
        "prevention": [
            "Early planting before peak moth season",
            "Crop rotation with non-host crops",
            "Use push-pull technology with Desmodium",
            "Monitor fields weekly during vegetative stage",
        ],
    },
    "Cassava___Mosaic_Disease": {
        "crop": "Cassava",
        "disease": "Cassava Mosaic Disease (CMD)",
        "severity_levels": ["mild", "moderate", "severe"],
        "treatments": {
            "organic": [
                "Remove and burn infected plants",
                "Plant CMD-resistant varieties (NASE 14, NAROCASS 1)",
                "Control whiteflies with neem spray",
            ],
            "chemical": [
                "Imidacloprid for whitefly control",
                "No cure - focus on prevention",
            ],
        },
        "prevention": [
            "Use certified disease-free planting material",
            "Plant resistant varieties",
            "Control whitefly populations",
            "Remove infected plants immediately",
        ],
    },
    "Tomato___Late_Blight": {
        "crop": "Tomato",
        "disease": "Late Blight",
        "severity_levels": ["early", "moderate", "advanced"],
        "treatments": {
            "organic": [
                "Copper-based fungicide spray",
                "Baking soda solution (1 tbsp per gallon)",
                "Remove infected leaves immediately",
            ],
            "chemical": [
                "Ridomil Gold (50g per 20L water)",
                "Mancozeb (50g per 20L water)",
                "Chlorothalonil spray",
            ],
        },
        "prevention": [
            "Avoid overhead irrigation",
            "Stake plants for air circulation",
            "Apply preventive fungicide before rains",
            "Rotate with non-solanaceous crops",
        ],
    },
}

def get_all_classes() -> list:
    """Get all disease classes for model training."""
    classes = set(PLANTVILLAGE_CLASSES)
    classes.update(CASSAVA_CLASSES)
    for crop_classes in UGANDA_CROP_CLASSES.values():
        classes.update(crop_classes)
    return sorted(list(classes))


def get_crop_from_class(class_name: str) -> str:
    """Extract crop name from disease class."""
    if "___" in class_name:
        return class_name.split("___")[0]
    return "Unknown"


def get_disease_from_class(class_name: str) -> str:
    """Extract disease name from disease class."""
    if "___" in class_name:
        return class_name.split("___")[1].replace("_", " ")
    return class_name
