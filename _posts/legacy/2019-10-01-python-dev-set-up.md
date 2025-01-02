---
title: "Setting up a Python project"
toc: true
toc_label: "Sections"
published: false
---

This guide is for Windows users, to help us set up Python projects systematically and consistently.

References:

- https://www.datacamp.com/community/tutorials/python-developer-set-up

First, install **Chocolatey**.

## Install Python

- Install **Python**. Good thing about Chocolatey is the automatic adding of python executable to PATH environment variable in Windows. (At the time of writing, latest Python version is 3.7.4)

  ```
  choco install python --version 3.7.4
  # or if you are upgrading
  choco upgrade python --version 3.7.4
  ```

- Upgrade **pip**.

  ```
  python -m pip install --upgrade pip
  ```

## Using VENV

I feel that using some form virtual environment, much like using containers, is the Python way of keep our dependencies isolated from other projects on the same machine, and to facilitate collaboration with other developers.

- Install **venv**. From experience, this library ships with Python3 unless you are using some special distro that had venv packaged separately from Python3.

  ```
  pip install -U venv
  ```

## (optional) Notebook

Popular for data science projects.

- Install Jupyter notebook.

  ```
  pip install -U jupyter
  ```

## Create Project

- Start a project folder, create a new virtual environment, and activate it.

  ```
  mkdir my-app
  cd my-app
  python -m venv new_venv

  # POSIX
  source new_venv/bin/activate
  # WINDOWS
  .\new_venv\Scripts\activate

  # deactivate environment
  deactivate
  ```

- Try installing a Python package in the virtual environment. It will only be stored in the current environment's Lib folder. A simple Python script to check where packages are installed in the current environment:

  ```
  >>> import sys
  >>> print(sys.prefix)
  ```

## (Optional) Cookiecutter Scaffolding

- Install **cookiecutter** to scaffold Python projects. Use a cookiecutter template to set up project directory structure. (TBC need to research more on uses of cookiecutter);

  ```
  pip install -U cookiecutter
  ```

## Managing Dependencies

- Save current dependencies of project.

  ```
  pip freeze > requirements.txt
  ```

- When working on this same project after cloning the repository, install all the dependencies again.

  ```
  pip install -r requirements.txt
  ```
