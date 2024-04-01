## Description

This CLI tool performs four main steps:
1. Reads the contents of the input Markdown file;
2. Uses some Go external libraries to parse Markdown and generate a valid
   HTML block;
3. Wraps the results with an HTML header and footer;
4. Saves the buffer to an HTML file that you can view in a browser.

## Supported flags
- `-file <file_name>`\
Markdown file to be previewed.
- `-s`\
Skip auto-preview.

## Usage
To build the app, run the following command in the root folder:

```
> go build .
```
Above command will generate `mdp` file. This name is defined in the `go.mod` file, and it will be the initialized module name.

After that you can run the program using the cmd and pass the file:

```
> .\todo-cli.exe -file File Name
```
