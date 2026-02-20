---
title: "Table Driven Tests in Python"
date: "2025-12-07"
categories: [Code]
tags: [programming]
math: false
layout: post
image: assets/table_driven_tests.jpg
---



Most, if not all, of the times, you know what any function in your program should do. You know what the inputs and outputs looks like. You know that even before you start thinking about writing the function. After writing the function, you check it with bunch of example inputs to see if it puts up to your expectations. Unit tests basically do the same. 

My approach to unit testing has remained more or less the same since I started writing them. But when I started learning go programming language, I came across this really simplistic, elegant and beautiful way to write unit test. They call it Table Driven testing, read [more here](https://go.dev/doc/tutorial/add-a-test). Read more about test in go at here. 

## Why Should We Write Tests?

There is no particular answer to why would you want to write tests? But there are far too many benefits to ignore writing tests. 

- Tests let's you check your code's intended behavior. 
- Tests acts as documentation for your code's intended behavior.
- Tests makes debugging easier by pinpointing the issue. 
- Tests makes collaboration easier by giving confidence in teammates code. 
- Tests make code reviews easier, if tests Fails, you can just reject PR. 

I can go on writing this list for the whole of this post but that's not the point of this pots. The point of this post is to 
## Simple Unit Tests

Let's say we are writing a function that _simplifies_ a string, meaning that it removes all the punctuation and trailing white spaces and converts it into lowercase. I would write a small python function to achieve that:

```python

def simplify(s: str) -> str:
    '''remove punctuations and make lowercase for a string'''
    
    punctuation = '''!"#$%&\'()*+,-./:;<=>?@[\\]^_`{|}~'''
    trans_map = s.maketrans({p:"" for p in punctuation})
    
    #remove punctuation
    s = s.translate(trans_map)
    
    # remove trailing whitespaces
    s = s.strip()
    
    # lowercase
    s = s.lower()

    return s
```

To test this simple function, I would write a simple unit test using python's `pytest` library. You can learn more about the pytest library [here](https://pytest.org/)

```python
import pytest
def test_simplify():

	inputs = "Boots the bear!"
	outputs = "boots the bear"
	function_output = simplify(input)
	
	assert function_output == outputs
```

I used to write test like that when I was kid. But that does the job. See I have used `inputs` as a variable even though it's singular because `input` in python is a keyword (`in` is also a keyword in python). While we can use `input` as a variable but if then try to use it a keyword, it might cause a problem. Another convention is to use `inputs_` with an underscore. But either way is fine. 

This works well for one input, but if wan to test for multiple cases, we can do scale it using a dictionary. 

```python

def test_simplify():

    test_cases = [
        {
            "inputs": "Boots the bear!",
            "want": "boots the bear",
        },
        {
            "inputs": "The wonderful bear, Boots ",
            "want": "the wonderful bear boots",
        },
        {
			"inputs": "",
			"want": "",
		},

		{
			"inputs": ".......",
			"want": "",
		},
    ]

    for test_case in test_cases:
        got = simplify(test_case["inputs"])
        assert got == test_case["want"], f"{test_case["want"]=}, {got=} "
```

 This is beautiful, with this we can test can test a few edge cases and determine if our code holds up to them before it crashes in the production level. We can run these test cases using `pytest` in our command line:

```bash
shyam@laptop: pytest
platform linux -- Python 3.13.5, pytest-8.3.4, pluggy-1.5.0
rootdir: /home/shyam/github/whateverproject
plugins: anyio-4.7.0
collected 4 items                                                                 

tests.py .....                                                   [100%]

==================== 1 passed in 0.01s =====================
```

Everything passes. And that's beautiful. But we can do a little more here. Sometimes we know that the test will fail in certain cases, we expect some error, and if our test is bypassing these errors then that might a bad sign. We want to detect that early. So we write a error prone test just to test it out. 

```python
def test_simplify():


    test_cases = [
        {
            "args": "Boots the bear!",
            "want": "boots the bear",
            "want_error": False
        },
        {
            "args": "The wonderful bear, Boots ",
            "want": "the wonderful bear boots",
            "want_error": False
        },
        {
            "args": 23,
            "kwargs": {},
            "want": "the wonderful bear, boots",
            "want_error": True
        },
    ]

    for test_case in test_cases:
        if test_case["want_error"]:
            flag = False
            try: 
                got = simplify(test_case["args"])
            except:
                flag = True
            assert flag , f"Wanted an error but got none"
            continue
        got = simplify(test_case["inputs"])
        assert got == test_case["want"], f"{test_case["want"]=}, {got=} "

```

Now we are talking. There we can see that if we get an integer instead of a string, the programs should fall apart and if it doesn't then something fishy is surely going on.

Now we are heading towards table driven testing. You can notice that the `test_cases` variable looks like a `json` file. It can be saved alone as json, and can also be made as a table. That's why we are calling it table driven tests. We can abstract a few thing out of here. 

You see that the bottom part of testing logic will basically be the same for all the functions we want to test. So I can write it in a separate function:

```python
import ast
def RUN(function: ast.FunctionDef , test_cases: list[dict]):

    for test_case in test_cases:

        if "kwargs" not in test_case.keys(): test_case["kwargs"]={}
        if "want_error" not in test_case.keys(): test_case["want_error"]=False

        if test_case["want_error"]:
            flag = False
            try: 
                got = function(*test_case["args"], **test_case["kwargs"])
            except:
                flag = True
            assert flag , f"Wanted an error but got none"
            return

        got = function(*test_case["args"], **test_case["kwargs"])
        assert got == test_case["want"], f"{test_case["want"]=}, {got=} "
```

This `RUN` will run test cases for all the functions for a bunch of testing pairs. And the logic even simplifies now:

```python
def test_simplify():


    test_cases = [
        {
            "args": ["Boots the bear!"],
            "want": "boots the bear",
            "want_error": False
        },
        {
            "args": ["The wonderful bear, Boots "],
            "want": "the wonderful bear boots",
            "want_error": False
        },
        {
            "args": [23],
            "kwargs": {},
            "want": "the wonderful bear, boots",
            "want_error": True
        },
    ]

    RUN(mod.simplify, test_cases)
```

Now this looks really simple and beautiful. 

>==NOTE==: `args` is not a simple string anymore, it's a list. (Well technically they should the tuple.) You might make that mistake so keep that in your mind.

This is not really a well written `RUN` function because it doesn't specify the kind of exception we are looking for it just looks for an error to pass the test. So if the error is because of a different reason, our test will still pass. That could be dangerous as well. So ideally we should include that Exception too in our tables.

```python
test_Case = {
            "args": ["Boots the bear!"],
            "want": "boots the bear",
            "want_error": False
            "exception": AttributeError
        },
```

And we should update our `RUN` function appropriately, we will do that shortly. Before that we need to address the lots of if blocks in there. 

Golang handles these test cases by making a struct of example cases. We don't have structs in python. The closest thing to a struct in python is a `dataclass`. We will define the test `dataclass` as the following:

```python
from dataclasses import dataclass, field
from typing import Any, Callable, Type, Tuple

@dataclass
class Case:
    #inputs
    args: Tuple[Any, ...]
    want: Any | Tuple[Any] | None = None

    kwargs: dict[str, Any] | None = field(default_factory=dict)
    exception: Type[Exception] | Tuple[Type[Exception], ...] = Exception
    want_error: bool = False
    name: str | None = None
```

At first glance, this might look a lot more complicated then our simple list of dictionaries but we aren't doing much here then defining the datatypes of the same variables. There are `args` as the list (tuple) of `Any`(which means literally any type). The want argument is an `Any` `|` (this pipe means `or`) tuple because there can be multiple outputs, and in Python, they are stored as tuples. (Here we have an edge over golang.) `kwargs` is called keyword arguments. They are like named arguments in a function. The optional type for that was suggest by GPT after a lot of debugging. 

The `exception` field is the most interesting one. Since Exception doesn't have a build in type in python, we have used `Type` to convert. It can also be a tuple. See the default value I have written is literal Exception. That was my intuition because when we have a particular exception while using the `try-except` block, we use that otherwise we do something akin to:

```python
try:
    f(1,0)
except Exception:
    print("Something")
```

Now we can simplify the run function:

```python
def RUN(function: Callable, test_cases: list[Case]):

    for test_case in test_cases:
        if test_case.want_error:
            with pytest.raises(test_case.exception):
                function(*test_case.args, **test_case.kwargs)
        else:
            got = function(*test_case.args, **test_case.kwargs)
            assert got == test_case.want, f"{test_case.want=}, {got=} "
```

That is very concise and beautiful. We have used `pytest.raises` context protocol. This will call the function inside the protocol. If the exception passed here is raised during execution, the test passes, if this exception isn't raised or some other exception is raised then the test fails. If we aren't expecting any error, we will just go and compare the outputs. Now we need to pass in the list of `Case` objects instead of dictionaries, so I have done that in a following way:

```python
def test_simplify():

    test_cases = [
        {
            "args": ["Boots the bear!"],
            "want": "boots the bear",
        },
        {
            "args": ["The wonderful bear, Boots "],
            "want": "the wonderful bear boots",
        },
        {
            "args": [23],
            "kwargs": {},
            "want": "the wonderful bear, boots",
            "want_error": True,
            "exception": AttributeError
        },
    ]

    test_cases = [Case(**things) for things in test_cases]
    RUN(mod.simplify, test_cases)
```

If we want to add a new test case, we just need to add a dictionary in the list. And we can keep appending this until we are satisfied. And there we have it, a beautiful, concise and _pythonic_ way to write table-driven tests. 

## Limitations 

- I have written this for functions with returning values. I haven't yet tried this approach for executive functions or _void_ functions. 
- This might be hard to integrate for the cases with dependency like databases, this is more helpful to test the helper functions. 

While this is a beautiful framework for writing tests when we have a good `wanted` and `got` pairs. (remember the Leetcode problems?) This isn't the only way to write test case. You can write them however you want. This is just an interesting approach inspired by golang. 

## When Should We Write Test Cases?

I have heard people say that we should write test cases for every function that we write. I don't agree with that. We should try to write test cases for every function but we should avoid some scenes as well. 

We shouldn't write test cases where the input depends on something external like a big file or a database. That should be tested separately in a manual or some other way. For example, my work involves reading a large file called `.lst` files and then extracting information. While I can collect a bunch of different files to write cases for them, I should avoid that. We might expose valuable company information while testing. We might have to gather a lot of data just to write the one test case which brings me to the second point. 

If setting up a test case is more expensive than the value it provides then we should shy away from it. And that happens when we have a lot of external dependencies. For example, to test a function that does some database stuff, we might need to install a database and set up the while thing. These cases should be broken down into simpler helper functions, and only those helper functions should be tested. 

If a helper function is to repeated many times is used widely in very different areas in the codebase, then we should definitely write test cases for it. But if the function is called just once, then we are wasting our time. And also if a function is very simple, e.g. does only lowercase in our example, then there is no point writing a test case in that. 

## Testing Automation with GitHub Actions 

Now that we have set up a few test cases, we should connect them with GitHub actions. There are following benefits:

- It runs automatically on every push, and if a test fails, it sends us a nice email. 
- It makes collaboration easier by checking on every PR. 
- It informs you if you refactoring breaks anything. So you can go back to previous comments. 

You can check python testing workflows in GitHub actions marketplace. I use this for all my repositories. Here is a test file.

```yml
name: RAG Testing

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

permissions:
  contents: read

jobs:
  build:

    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4
    - name: Set up Python 3.13
      uses: actions/setup-python@v3
      with:
        python-version: "3.13"
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install flake8 pytest
        pip install -r requirements.txt  
        if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
    - name: Lint with flake8
      run: |
        # stop the build if there are Python syntax errors or undefined names
        flake8 . --count --select=E9,F63,F7,F82 --show-source --statistics
        flake8 . --count --exit-zero --max-complexity=10 --max-line-length=127
    - name: Test with pytest
      run: |
        pytest
``` 

I haven't written this, I have directly copied from GitHub marketplace. but you can copy and paste this too. To use this, make a `.github/workflows` directory in your root, and put this code in a called called `actions.yml`. This should work just as fine. Make sure you have a `requirements.txt` file.
## Conclusion

Writing good tests has always been an industry standard in the tech community. Most, if not all, the companies in the world implement unit test cases to make software developments standard and smooth. But outside the corporate test cases are rather underrated. I think there are many interesting uses of test cases.
If you are a computer science teacher, you an use test cases to grad assignments from your students. You can also put it through PR and use GitHub actions to auto grad the assignments. 
If you are a researcher, and you are trying different models (for moving the droplet in chemical field for example.) You can write test cases to determine the level of physics your model is reaching. 
If you are writing blog posts like this, you should also implement test cases to prevent yourself from publishing bullhit online. 

Anyway I hope you have enjoyed this blog post. If you did, do read, share it. And if you want to give a feedback you can reach me at [my email](mailto:shyam10kwd@gmail.com).
