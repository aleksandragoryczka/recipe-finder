# recipe-finder

## Table of contents

- [Introduction](#introduction)
- [Technologies used](#technologies-used)
- [Launching](#launching)
- [Usage example](#usage-example)


## Introduction
This project helps users to generate list of meals based on their fridge contents with minimal number of missing 
ingredients needed.

## Technologies used
- Go
- PostgreSQL

## Launching
In main folder execute command o build program: `make build`

Then run program with your flags:

    --ingredients - list of ingredients from your fridge
    --numberOfRecipes - the maximum number if recipes you want

Example run: `./recipeFinder --ingredients=carrots,eggs --numberOfRecipes=4`

If you want more details about program and its running: `./recipeFinder -h`.




