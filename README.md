# datagen
A simple data generation script to generate files and folders

# Features
- Generates a base folder with a predefined folder structure based on real-world use cases.
- Minimum file count is customizable via command-line flags.
- Creates files with some initial data and performs random operations including adding new files, modifying existing files, and removing some data from existing files.
- Capable of adding and removing folders randomly.
- Compatible with both Linux and Windows operating systems.

# Build
```
go build -o datagen main.go
```

# Command Line Flags & Usage
`--base`: Base folder for the test data (Default: "testData")
`--minFiles`: Minimum number of files to be created in each folder (Default: 5)
`--randomOps`: Perform random operations on the data (Default: true)

```
./datagen -base myData -minFiles 3 -randomOps true
```
# Notes
- Ensure that your file and folder paths are within the limits of your operating system to avoid any issues.
- Feel free to contribute to the improvement of this tool by creating Pull Requests or reporting Issues.
