export const getMinioURL = (path: string) => {
  return `${import.meta.env.VITE_MINIO_URL}${path}.png`
}