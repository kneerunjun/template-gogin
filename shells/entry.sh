#! /bin/sh
# incase the command line arguments arent as expected this will emit a message on the command line
usage() { echo "Usage: $0 [-v <true/false>] [-f <true/false>] [-s <true/false>]" 1>&2; exit 1; }
_term(){
    # incase of interruption this helps in gracefully shutting down the container 
    # App code will have to handle the kill message that we send in here
    echo "shutting down the application container"
    # /usr/sbin/crond stop
    kill -TERM "$child" 2>/dev/null
}

trap _term SIGTERM #so as to pass it down
echo "starting cron deamon..."
# /usr/sbin/crond -f -l 8&


# getting all the command line arguments 
verbose="false"
filelog="false"
while getopts ":v:f:s:" o; do
    case "${o}" in
        v)
            verbose=${OPTARG}
            ;;
        f)
            filelog=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done
echo $verbose
echo $filelog
# TODO: change the app name here else the entrypoint to the container will fail
/usr/bin/template-app -verbose $verbose -flog $filelog&

# waiting for the go application
child=$!
wait "$child"
# enable this below line incase you are using crons from within the container 
# /usr/sbin/crond stop # if the go app container exits gracefully without any user interruption 