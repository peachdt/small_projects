# Utility Tool
Purpose of this project is to make rebuilding magento vagrant box a lot easier by automating all the removing/creating tables, setting up files in onestep.

# Usage
Since this is a personal project for internal company use, many of actual ansible setup files with creds are not pushed. This repo is just for personal project tracking.
Run `go run main.go` for more details on the arguments required (-action=rebuild -target_box=vagrant_box_name -env=local_box or server_box )
Depending on the env, different ansible template files will be used to write to the actual devops/ansible repo and replace the `replace_me` with target_box names.
db folder is used to connect to remote mysql to delete/rebuild tables used for Magento.