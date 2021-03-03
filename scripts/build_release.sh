#!/bin/bash

verbose=false
artifact_dir="artifacts"
gh_api="https://api.github.com"
gh_owner="reecerussell"
gh_repo_name="migrations"
gh_repo="$gh_api/repos/$gh_owner/$gh_repo_name/releases?access_token=$GITHUB_TOKEN"

exit_if_error() {
    if [[ $1 -ne 0 ]]; then
        log_info "\nERROR!"

        if [[ $2 != "" ]]; then
            log_info "> Message: $2"
        fi

        if [[ $3 != "" ]]; then
            log_info "> Details: $3"
        fi

        exit $1
    fi
}

log_info() {
    echo -e $1
}

log_verbose() {
    if [[ $verbose == true ]]; then
        echo $1
    fi
}

build_amd64() {
    os=linux
    output="$artifact_dir/migrations"

    os_pat="^(windows|linux)$"
    if [[ $1 =~ $os_pat ]]; then
        os=$1

        if [[ $1 == "windows" ]]; then
            output="$output.exe"
        fi
    fi

    log_info "\n---"
    log_info "Build for $os (amd64)"
    log_verbose "> Output: $output"
    log_info "---\n"

    main=cmd/main.go
    log_info "Building $main..."
    GOOS=$os GOARCH=amd64 CGO_ENABLED=0 go build -o $output $main &> out.txt
    exit_if_error $? "An error occurred while building" "$(cat out.txt && rm out.txt)"

    log_info "Built successfully!"
}

download_modules() {
    log_info "\n---"
    log_info "Modules"
    log_info "---\n"

    log_info "Downloading..."
    go mod download &> out.txt
    exit_if_error $? "An error occurred while downloading modules" "$(cat out.txt && rm out.txt)"

    log_info "Verifying modules..."
    go mod verify > out.txt
    exit_if_error $? "An error occurred while verifying modules" "$(cat out.txt)"

    out=$(cat out.txt && rm out.txt)
    if [[ $(echo $out | wc -l) -gt 0 ]]; then
        log_verbose "$out"
    fi
}

create_release() {
    latest_tag=$(git tag --sort=committerdate -l | tail -1)
    latest_commit_message=$(git log -1 --pretty=format:"%s")
    log_info "Creating release for $latest_tag, with message '$latest_commit_message'..."

    release_json=$(printf '{"tag_name": "%s","target_commitish": "master","name": "%s","body": "%s","draft": false,"prerelease": false}' \
        "$latest_tag" "$latest_tag" "$latest_commit_message")
    log_verbose "Release JSON: $release_json"

    log_verbose "Making release request..."
    curl --data "$release_json" $gh_repo > out.txt
    exit_if_error $? "An error occurred while creating the release" "$(cat out.txt)"

    release_id=$(cat out.txt | jq '.id')
    if [[ "$release_id" == "null" ]]; then
        exit_if_error 1 "Failed to create release" "$(cat out.txt && rm out.txt)"
    fi

    rm out.txt
    log_verbose "Created release $release_id"

    log_info "Uploading artifacts..."

    for f in "$artifact_dir"/*; do
        if [[ -d "$f" ]]; then
            continue
        fi

        log_info "Uploading '$f'"

        asset="https://uploads.github.com/repos/$gh_owner/$gh_repo_name/releases/$release_id/assets?name=$(basename $f)"
        log_verbose "Uploading to: $asset"

        log_verbose "Posting artifact to Github..."
        curl --data-binary @"$f" \
            -X POST \
            -H "Accept: application/vnd.github.v3+json" \
            -H "Content-Type: application/octet-stream" \
            -H "Authorization: token $GITHUB_TOKEN" \
            "$asset" &> out.txt
        exit_if_error $? "An error occurred while uploading '$f'" "$(cat out.txt | tail -1 && rm out.txt)"

        log_verbose "Successfully posted artifact"
    done

    log_info "Release successful!"
}

main() {
    if [[ "$VERBOSE" == "1" ]]; then
        verbose=true
    fi

    version=$(git tag --sort=committerdate -l | tail -1)

    log_info "Building Release!"
    log_info "> VERSION: $version"
    log_info "> Verbose logging: $verbose"
    log_verbose "> Git repository: $gh_repo_name"


    download_modules

    build_amd64 linux
    build_amd64 windows

    create_release

    exit 0
}

main