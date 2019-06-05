import Uppy from '@uppy/core'
import Dashboard from '@uppy/dashboard'
import AwsS3 from '@uppy/aws-s3'

var uppy = Uppy({
    debug: true,
    autoProceed: true,
    restrictions: {
        maxFileSize: 1000000,
        // maxNumberOfFiles: 2,
        minNumberOfFiles: 1,
        allowedFileTypes: ['image/*']
    },
    showProgressDetails: true,
    note: 'Images only, 2–3 files, up to 1 MB'
})

uppy.use(Dashboard, {
    inline: true,
    target: '#drag-drop-area',
    metaFields: ['name','key'],
    width: '100%',
    height: '100%',
    // thumbnailWidth: 280,
    // defaultTabIcon: defaultTabIcon,
    // showLinkToFileUploadResult: true,
    // showProgressDetails: false,
    // hideUploadButton: false,
    // hideRetryButton: false,
    // hidePauseResumeButton: false,
    // hideCancelButton: false,
    // hideProgressAfterFinish: false,
    // note: null,
    // closeModalOnClickOutside: false,
    // closeAfterFinish: false,
    // disableStatusBar: false,
    // disableInformer: false,
    // disableThumbnailGenerator: false,
    // disablePageScrollWhenModalOpen: true,
    // animateOpenClose: true,
    // proudlyDisplayPoweredByUppy: true,
    // onRequestCloseModal: () => this.closeModal(),
    // showSelectedFiles: true,
    // locale: defaultLocale,
    browserBackButtonClose: false
})
uppy.use(AwsS3, {
    getUploadParameters(file) {
        console.log(file)
        return fetch('/sign', {
            method: 'post',
            // Send and receive JSON.
            headers: {
                accept: 'application/json',
                'content-type': 'application/json'
            },
            body: JSON.stringify({
                filename: file.name,
                contentType: file.type,
                size: file.size
            })
        }).then((response) => {

            return response.json()
        }).then((data) => {
            file.meta["bucket"] = data.bucket;
            file.meta["key"] = data.key;
            // Return an object in the correct shape.
            return {
                method: data.method,
                url: data.url,
                fields: data.fields
            }
        })
    }
})
uppy.on('upload-success', (file, data) => {
    console.log('upload success')
    console.log('file:',file)
    console.log('data:',data)
})
uppy.on('complete', (result) => {
    console.log(result);
    document.getElementById('FileName').value = result.successful[0].name;
    document.getElementById('FilePath').value = result.successful[0].uploadURL;
    document.getElementById('FileSize').value = result.successful[0].size;
    document.getElementById('FileExtension').value = result.successful[0].extension;
    document.getElementById('FileKey').value = result.successful[0].meta.key;
    document.getElementById('FileBucket').value = result.successful[0].meta.bucket;
});
