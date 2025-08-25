#!/bin/bash
# ./scripts/publish.sh <zipfile> <params>
# ./scripts/publish.sh 
# ./scripts/publish.sh hello-go.zip
# ./scripts/publish.sh hello-go --delete
# ./scripts/publish.sh --all                                # publish all lambdas

if [ "$#" -lt "1" ]; then
    echo "Specify function_name as first argument"
    exit 1
fi

endpoint="http://localhost:4566"
zipfile=$1                          # the zip file should has the same name of the function
function_name=${zipfile%.zip}
binpath=bin # root of workspace
bin_lambda_path=$binpath/lambdas
zipfile_lambda="fileb://$bin_lambda_path"

# cd bin || exit

runtime=go1.x
timeout=150
handler="bootstrap"

# if flag -delete, remove

echo "=== Debug variables ==="
echo "zipfile: $zipfile"
echo "function_name: $function_name"

lambda_exists() {
    local function_name=$1
    # aws --profile localstack --endpoint-url=http://localhost:4566 lambda list-functions
    # aws --profile localstack --endpoint-url=http://localhost:4566 lambda get-function --function-name $function_name
    # aws --profile localstack --endpoint-url=http://localhost:4566 lambda get-function --function-name hello-go

    output=$(aws --profile localstack --endpoint-url=$endpoint lambda get-function --function-name "$function_name" 2>&1)
    error=$?

    if [ $error -eq 0 ]; then
        echo "Lambda function '$function_name' exists."
    elif echo "$output" | grep -q 'ResourceNotFoundException'; then
        echo "Lambda function '$function_name' does not exist."
    else
        echo "Error checking function: $output"
    fi

    return $error
}

create_lambda(){
    local function_name=$1
    local runtime=$2
    local handler=$3
    local timeout=$4
    local zipfile=$5

    echo "creating lambda $function_name $runtime $handler $timeout $zipfile"

    output=$(
        aws --profile localstack --endpoint-url=$endpoint \
        lambda create-function \
        --function-name "${function_name}" \
        --runtime "${runtime}" \
        --role arn:aws:iam::000000000000:role/lambda-role \
        --handler "${handler}" \
        --timeout "${timeout}" \
        --zip-file "${zipfile}"
    )
    error=$?
    if [ $error -ne 0 ]; then
        echo "Error create function: $output"
    fi
}

update_lambda(){
    local function_name=$1
    local runtime=$2
    local handler=$3
    local timeout=$4
    local zipfile=$5

    echo "updating lambda $function_name $runtime $handler $timeout $zipfile"

    output=$(
        aws --profile localstack --endpoint-url=$endpoint \
        lambda update-function-code \
        --function-name "${function_name}" \
        --zip-file "${zipfile}"
    )
    error=$?
    if [ $error -ne 0 ]; then
        echo "Error updating function: $output"
    fi
}

delete_lambda(){
    aws --profile localstack --endpoint-url=$endpoint \
    lambda delete-function \
    --function-name "${function_name}"
}

### main ###----------------------
DELETE=false
ALL=false

for arg in "$@"; do
  if [ "$arg" = "--delete" ]; then
    DELETE=true
    break
  fi
  if [ "$arg" = "--all" ]; then
    ALL=true
    break
  fi
done

if [ "$DELETE" = true ]; then
    echo "deleting lambda $function_name"
    delete_lambda "$function_name"
    exit 0 
fi

if [ "$ALL" = true ]; then
    # publish all files of bin/lambdas folder
    echo "processing all ..."
    for file in "$bin_lambda_path"/*.zip; do
        # Get the base filename (e.g., hello-go.zip)
        zipfile="${file##*/}"
        function_name="${zipfile%.zip}" # remove extension zip
        
        if ! lambda_exists "$function_name"; then
            create_lambda "$function_name" $runtime $handler $timeout "$zipfile_lambda/$zipfile"
        else
            update_lambda "$function_name" $runtime $handler $timeout "$zipfile_lambda/$zipfile"
        fi
    done
else
    # publish one file by time
    if lambda_exists "$function_name"; then
        update_lambda "$function_name" $runtime $handler $timeout "$zipfile_lambda/$zipfile"
    else
        create_lambda "$function_name" $runtime $handler $timeout "$zipfile_lambda/$zipfile"
    fi
fi

echo "finished"