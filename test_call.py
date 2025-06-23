import inspect

def caller_function():
    return (inspect.stack()[1].function, inspect.stack()[1].filename)

def first_function():
    print(f"I am 'first_function'. My caller is: {caller_function()}")
    second_function()

def second_function():
    print(f"I am 'second_function'. My caller is: {caller_function()}")

if __name__ == "__main__":
    first_function()
