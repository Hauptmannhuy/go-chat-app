import { useState } from "react"
import { FileImage } from "../types/image"


export function useUploadImage() {

  const [inputImage, setImage] = useState(new FileImage())

/**
* @param {File} file 
*/
  const uploadFile = (file) => {
    const reader = new FileReader()
    
    reader.onload = function () {
      const bytes = new Uint8Array(this.result);

      window.createImageBitmap(file).then((bitmap) => {
        const newFileImage = new FileImage(bytes, file.type, bitmap.width,bitmap.height);
        setImage(newFileImage)
      });
    }

    reader.readAsArrayBuffer(file)
  }


  return {inputImage, uploadFile}

}