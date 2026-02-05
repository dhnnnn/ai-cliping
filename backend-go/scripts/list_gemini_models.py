#!/usr/bin/env python3
"""
Script to list available Gemini models for your API key
"""

import google.generativeai as genai
import os
import sys

def list_models(api_key):
    """List all available Gemini models that support generateContent"""
    
    # Configure API
    genai.configure(api_key=api_key)
    
    print("=" * 60)
    print("GEMINI AVAILABLE MODELS")
    print("=" * 60)
    print()
    print("Models that support generateContent:")
    print("-" * 60)
    
    models_found = []
    
    try:
        for model in genai.list_models():
            if 'generateContent' in model.supported_generation_methods:
                models_found.append(model.name)
                print(f"✓ {model.name}")
                print(f"  Display name: {model.display_name}")
                print(f"  Description: {model.description}")
                print()
        
        print("=" * 60)
        print(f"Total: {len(models_found)} models available")
        print("=" * 60)
        
        if models_found:
            print("\nRECOMMENDATIONS:")
            print("- For fast, free: Look for models with 'flash' in the name")
            print("- For best quality: Look for models with 'pro' in the name")
            print("\nUpdate your Go code to use one of these model names.")
        else:
            print("\n⚠️  No models found! Check your API key.")
            
    except Exception as e:
        print(f"\n❌ ERROR: {e}")
        print("\nPossible issues:")
        print("- Invalid API key")
        print("- Network connection problem")
        print("- API quota exceeded")
        sys.exit(1)

if __name__ == "__main__":
    # Try to get API key from environment or command line
    api_key = os.getenv("GEMINI_API_KEY")
    
    if not api_key and len(sys.argv) > 1:
        api_key = sys.argv[1]
    
    if not api_key:
        print("Usage:")
        print("  python list_gemini_models.py YOUR_API_KEY")
        print("Or set GEMINI_API_KEY environment variable")
        print("\nExample:")
        print("  python list_gemini_models.py AIzaSy...")
        sys.exit(1)
    
    list_models(api_key)
