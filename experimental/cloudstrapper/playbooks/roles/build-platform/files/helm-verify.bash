
gitPublishBranch=$3
gitUserName=$1
gitPat=$2
#helm repo add *requires* PAT and will not work with password


helm repo add $gitPublishBranch --username $gitUserName --password $gitPat https://raw.githubusercontent.com/$gitUserName/$gitPublishBranch/master/
helm repo update && helm repo list
helm search repo $gitPublishBranch

