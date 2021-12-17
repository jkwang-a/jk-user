#!/bin/bash

\cp ../bmc-user user/user/

tag="v`date +%Y%m%d%H%M`"

work_dir=$(cd "$(dirname $0)";pwd)


for project in `find ./ -type d -d 1 | awk -F'/' '{print $NF}'`
do
	cd $work_dir/$project; docker build -t reg.docker.zenlayer.net/bmc-admin-test/${project}:$tag .	
	docker push reg.docker.zenlayer.net/bmc-admin-test/${project}:$tag
done


echo TAG is "$tag"
