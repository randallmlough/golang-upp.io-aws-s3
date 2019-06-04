import Uppy from '@uppy/core'
import Dashboard from '@uppy/dashboard'
import AwsS3 from '@uppy/aws-s3'

var uppy = Uppy({
    debug: true
})

uppy.use(Dashboard, {
    inline: true,
    target: '#drag-drop-area'
})
uppy.use(AwsS3, {
    getUploadParameters(file) {
        return fetch('/sign', {
            method: 'post',
            // Send and receive JSON.
            headers: {
                accept: 'application/json',
                'content-type': 'application/json'
            },
            body: JSON.stringify({
                filename: file.name,
                contentType: file.type
            })
        }).then((response) => {
            // Parse the JSON response.
            return response.json()
        }).then((data) => {
            // Return an object in the correct shape.
            return {
                method: data.method,
                url: data.url,
                fields: data.fields
            }
        })
    }
})
uppy.on('complete', (result) => {
    console.log(result);
    console.log(result.successful[0].name);
    console.log(result.successful[0].uploadURL);
    console.log('Upload complete! We’ve uploaded these files:', result.successful);
});
