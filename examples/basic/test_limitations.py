# Test script to demonstrate limitations and alternatives in embedded Python

def test_threading_alternative():
    """Demonstrate that threading works as an alternative to multiprocessing"""
    import threading
    import time
    import queue
    
    print("Testing threading alternative to multiprocessing...")
    
    def worker(work_id, result_queue):
        # Simulate some work
        result = work_id * work_id
        time.sleep(0.1)  # Simulate processing time
        result_queue.put(f"Worker {work_id}: {result}")
    
    # Create queue for results
    result_queue = queue.Queue()
    threads = []
    
    # Start worker threads
    for i in range(5):
        t = threading.Thread(target=worker, args=(i, result_queue))
        t.start()
        threads.append(t)
    
    # Wait for all threads to complete
    for t in threads:
        t.join()
    
    # Collect results
    results = []
    while not result_queue.empty():
        results.append(result_queue.get())
    
    return results

def test_concurrent_futures():
    """Demonstrate concurrent.futures as multiprocessing alternative"""
    from concurrent.futures import ThreadPoolExecutor
    import time
    
    print("Testing concurrent.futures alternative...")
    
    def cpu_bound_task(n):
        # Simulate CPU-bound work
        result = 0
        for i in range(n * 100000):
            result += i
        return f"Task {n}: {result}"
    
    # Use ThreadPoolExecutor instead of ProcessPoolExecutor
    with ThreadPoolExecutor(max_workers=3) as executor:
        futures = [executor.submit(cpu_bound_task, i) for i in range(5)]
        results = [future.result() for future in futures]
    
    return results

def test_numpy_compatibility():
    """Test that NumPy works well in embedded environment"""
    try:
        import numpy as np
        
        print("Testing NumPy compatibility...")
        
        # Create test data
        data = np.random.rand(1000, 1000)
        
        # Perform computation
        result = np.linalg.norm(data)
        
        return f"NumPy test successful: matrix norm = {result:.2f}"
    except ImportError:
        return "NumPy not available (install with: pip install numpy)"

def test_multiprocessing_failure():
    """Demonstrate that multiprocessing fails in embedded environment"""
    try:
        import multiprocessing
        
        def worker_func(x):
            return x * x
        
        print("Testing multiprocessing (expected to fail)...")
        
        # This will likely fail in embedded Python
        with multiprocessing.Pool(processes=2) as pool:
            results = pool.map(worker_func, [1, 2, 3, 4])
        
        return f"Multiprocessing unexpectedly succeeded: {results}"
        
    except Exception as e:
        return f"Multiprocessing failed as expected: {type(e).__name__}: {e}"

def get_environment_info():
    """Get information about the Python environment"""
    import sys
    import os
    
    info = {
        "python_version": sys.version,
        "platform": sys.platform,
        "executable": sys.executable,
        "argv": sys.argv,
        "path": sys.path[:3],  # First 3 entries
        "modules": len(sys.modules),
    }
    
    return info

def run_all_tests():
    """Run all compatibility tests"""
    print("=== Embedded Python Compatibility Tests ===\n")
    
    # Environment info
    print("1. Environment Information:")
    env_info = get_environment_info()
    for key, value in env_info.items():
        print(f"   {key}: {value}")
    print()
    
    # Threading test
    print("2. Threading Alternative Test:")
    threading_results = test_threading_alternative()
    for result in threading_results:
        print(f"   {result}")
    print()
    
    # Concurrent futures test
    print("3. Concurrent Futures Test:")
    future_results = test_concurrent_futures()
    for result in future_results:
        print(f"   {result}")
    print()
    
    # NumPy test
    print("4. NumPy Compatibility Test:")
    numpy_result = test_numpy_compatibility()
    print(f"   {numpy_result}")
    print()
    
    # Multiprocessing test (expected to fail)
    print("5. Multiprocessing Test (Expected to Fail):")
    mp_result = test_multiprocessing_failure()
    print(f"   {mp_result}")
    print()
    
    print("=== Test Complete ===")
    
    return {
        "environment": env_info,
        "threading_works": len(threading_results) > 0,
        "concurrent_futures_works": len(future_results) > 0,
        "numpy_available": "NumPy test successful" in test_numpy_compatibility(),
        "multiprocessing_fails": "failed as expected" in mp_result
    }

# For easy testing from Go
if __name__ == "__main__":
    run_all_tests()