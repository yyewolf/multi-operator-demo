# Instrumentation

The control panel is a simple web application that allows you to interact with the operators and see the state of the cluster. It is not part of the demo and is not required to run it. It is only used to help with the presentation.

It uses a bit of trickery to simplify the lifecycle of the demo. The App operator does not handle config refreshing to simplify things. The control panel can do that by doing the following:
- Annotate pods to trigger a configmap refresh
- Exec into the pod to tell supervisor to reread and update the config

The first problem can be solved by adding a step in the App operator.  
The second one can be harder to solve since supervisor does not have a way to reload the config via HTTP. One might choose to use another process manager that has this feature.

## Requirements

To run this control panel, you need to have the following tools installed:
- [Python 3](https://www.python.org/downloads/)

## Installation

Follow these steps to install the dependencies:

```bash
$ python3 -m venv .venv
$ source .venv/bin/activate
$ pip install -r requirements.txt
```

## Running the control panel

To run the control panel, in your terminal, run the following command:

```bash
$ source .venv/bin/activate
$ python app.py
```
