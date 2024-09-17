#arg1: gitusername arg2: gitpassword arg3: git publish branch 
gitBranchCheckout=a7580153
gitPublishBranch=$3
gitUserName=$1
gitPat=$2
dirMagmaRoot=$4
dirHelmChart=~/magma-charts

cd $dirMagmaRoot
git checkout $gitBranchCheckout
cd $dirHelmChart
git init
helm package $dirMagmaRoot/orc8r/cloud/helm/orc8r/ && helm repo index .
git add . && git commit -m "Initial Commit"
git remote add origin https://$gitUserName:$gitPat@github.com/$gitUserName/$gitPublishBranch && git push -u origin master

