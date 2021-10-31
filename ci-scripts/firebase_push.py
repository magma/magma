import pyrebase
import json
import time
import os

def main():
    # Read db config
    FIREBASE_CONFIG = os.environ['FIREBASE_CONFIG']
    config = json.loads(FIREBASE_CONFIG)
    # Initialize pyrebase app
    firebase = pyrebase.initialize_app(config)
    # Get a reference to the auth service
    auth = firebase.auth()
    # Log the user in
    user = auth.sign_in_with_email_and_password(config['auth_email'], config['auth_password'])
    # Get a reference to the database service
    db = firebase.database()

    # Grab environment variables
    workers_env = os.environ['WORKERS']
    build_info_env = os.environ['BUILD_INFO']
    workload_env = os.environ['WORKLOAD']
    build_id = os.environ['BUILD_ID']

    # Process env variables
    workers = [x.strip() for x in workers_env.split(',')]
    try:
        build_info = json.loads(build_info_env)
    except ValueError:  # includes simplejson.decoder.JSONDecodeError
        print('Decoding build_info JSON has failed: ', build_info_env)
        raise ValueError('Decoding build_info JSON has failed')

    try:
        workload = json.loads(workload_env)
    except ValueError:  # includes simplejson.decoder.JSONDecodeError
        print('Decoding workload JSON has failed: ', workload_env)
        raise ValueError('Decoding workload JSON has failed')

    # Push build information
    # Grab workload from environment
    print('Pushing the following build info to database: builds/{}'.format(build_id))
    print(build_info)
    db.child("builds").child(build_id).set(build_info, user['idToken'])

    # Push workloads
    # Grab workload from environment
    print('Pushing the following workload to workers: ', workers)
    print(workload)
    # Push new workload to all active workers. This generates a random ID for the workload
    for worker in workers:
        db.child("workers").child(worker).child("workloads").push(workload, user['idToken'])

if __name__ == "__main__":
    main()
