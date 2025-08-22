# Example Python module for testing Phase 3 function calling
import math

def calculate_area(radius):
    """Calculate the area of a circle given its radius."""
    return math.pi * radius ** 2

def factorial(n):
    """Calculate factorial of n."""
    if n <= 1:
        return 1
    return n * factorial(n - 1)

def statistics(numbers):
    """Calculate basic statistics for a list of numbers."""
    if not numbers:
        return {"error": "Empty list"}
    
    total = sum(numbers)
    count = len(numbers)
    mean = total / count
    
    sorted_nums = sorted(numbers)
    median = sorted_nums[count // 2] if count % 2 == 1 else (sorted_nums[count//2 - 1] + sorted_nums[count//2]) / 2
    
    return {
        "count": count,
        "sum": total,
        "mean": mean,
        "median": median,
        "min": min(numbers),
        "max": max(numbers)
    }

def process_data(data):
    """Process a dictionary of data and return enhanced results."""
    results = {
        "original": data,
        "processed_at": "2024",
        "summary": {}
    }
    
    if "numbers" in data:
        results["summary"]["number_stats"] = statistics(data["numbers"])
    
    if "text" in data:
        results["summary"]["text_length"] = len(data["text"])
        results["summary"]["word_count"] = len(data["text"].split())
    
    return results