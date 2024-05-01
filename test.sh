./build.sh

cd testdata/unstable/
go build

GOSTABLE_OUT=$(../../gostable . 2>&1)

# Initialize the CLIPPED_OUT variable
CLIPPED_OUT=""

# Process the output line by line
while IFS= read -r line; do
    # Remove the path before "gostable/testdata"
    modified_line="${line#*/gostable/testdata}"
    
    # Prepend "gostable/testdata" to the modified line
    modified_line="gostable/testdata$modified_line"
    
    # Append the modified line to CLIPPED_OUT
    CLIPPED_OUT+="$modified_line"$'\n'
done < <(echo "$GOSTABLE_OUT")

# Remove the trailing newline character from CLIPPED_OUT
CLIPPED_OUT="${CLIPPED_OUT%$'\n'}"

#echo "Modified output:"
#echo "$CLIPPED_OUT"

# Compare CLIPPED_OUT with the contents of the file
diff -u <(echo "$CLIPPED_OUT") golden

# Check the exit status of the diff command
if [ $? -eq 0 ]; then
    echo "gostable output matches testdata/unstable/golden"
else
    echo "gostable output does not match testdata/unstable/golden"
fi

cd ../stable
go build

GOSTABLE_OUT=$(../../gostable . 2>&1)

# Initialize the CLIPPED_OUT variable
CLIPPED_OUT=""

# Process the output line by line
while IFS= read -r line; do
    # Remove the path before "gostable/testdata"
    modified_line="${line#*/gostable/testdata}"
    
    # Prepend "gostable/testdata" to the modified line
    modified_line="gostable/testdata$modified_line"
    
    # Append the modified line to CLIPPED_OUT
    CLIPPED_OUT+="$modified_line"$'\n'
done < <(echo "$GOSTABLE_OUT")

# Remove the trailing newline character from CLIPPED_OUT
CLIPPED_OUT="${CLIPPED_OUT%$'\n'}"

#echo "Modified output:"
#echo "$CLIPPED_OUT"

# Compare CLIPPED_OUT with the contents of the file
diff -u <(echo "$CLIPPED_OUT") golden

# Check the exit status of the diff command
if [ $? -eq 0 ]; then
    echo "gostable output matches testdata/stable/golden"
else
    echo "gostable output does not match testdata/stable/golden"
fi


cd ..
