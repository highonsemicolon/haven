import boto3
import os
import zipfile
from io import BytesIO

DESTINATION_BUCKET_NAME = os.environ['DESTINATION_BUCKET']

s3_client = boto3.client('s3')

def lambda_handler(event, context):
    for record in event['Records']:
        SOURCE_BUCKET_NAME = record['s3']['bucket']['name']
        object_key = record['s3']['object']['key']
       

        response = s3_client.get_object(Bucket=SOURCE_BUCKET_NAME, Key=object_key)
        zip_data = response['Body'].read()
        
        # Unzip the file
        with zipfile.ZipFile(BytesIO(zip_data)) as zip_ref:
            for file_name in zip_ref.namelist():
                file_data = zip_ref.read(file_name)
                destination_key = file_name.replace('build', object_key, 1)

                # Upload each unzipped file to destination bucket
                s3_client.put_object(Body=file_data, Bucket=DESTINATION_BUCKET_NAME, Key=destination_key)
        
        return {
            'statusCode': 200,
            'body': 'Unzipping completed successfully!'
        }
