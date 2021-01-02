import axios from 'axios';
const baseURL = 'http://localhost:8090/';

export let uploadPhotos = (files: File[]) => {
  // we upload the photos individually so if one upload fails, the earlier uploads will still go through
  files.forEach((file) => {
    let formData = new FormData();
    formData.append('file', file, file.name);
    axios.post(baseURL + 'upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  });
};
