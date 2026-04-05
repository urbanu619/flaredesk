export const downLoadFile = (res) => {
  const {fileName, fileUrl } = res
  console.log(fileName, fileUrl)
  const link = document.createElement('a');
  link.href = fileUrl;
  link.target = '_blank' // 新标签页打开
  link.download = fileName + ".csv" ;
  document.body.appendChild(link)
  link.click();
  document.body.removeChild(link)
}
