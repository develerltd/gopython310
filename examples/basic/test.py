# Example Python file for testing Phase 2
print("Hello from test.py!")

def calculate_fibonacci(n):
    if n <= 1:
        return n
    return calculate_fibonacci(n-1) + calculate_fibonacci(n-2)

# Calculate and print some fibonacci numbers
for i in range(8):
    fib = calculate_fibonacci(i)
    print(f"fibonacci({i}) = {fib}")

# Test some Python features
data = [1, 2, 3, 4, 5]
squared = [x**2 for x in data]
print(f"Original: {data}")
print(f"Squared: {squared}")

# Test dictionary
person = {"name": "Alice", "age": 30, "city": "New York"}
print(f"Person: {person}")

print("File execution completed successfully!")