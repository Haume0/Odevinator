import { useProgress } from "./Store";

export function LoadingModal() {
  const [progress] = useProgress()
  return(
    <>
      {progress.state && (
        <div class=" fixed bg-black/60 backdrop-blur-sm w-screen flex items-center justify-center h-[100svh] inset-0">
          <section class="w-[32rem] p-12 bg-white items-center justify-center rounded-xl flex flex-col gap-4">
            {progress.value <=0 ? (
              <img src="/spinner.gif" class="size-24" alt="" />
            ):(
              <div class="relative w-full h-4 bg-gray-200 rounded-full">
                <div class="absolute h-full bg-blue-500 rounded-full" style={{width: `${progress.value}%`}}></div>
              </div>
            )}
            <h1 class="text-3xl font-bold">Lütfen bekleyin...</h1>
            <p class="font-medium">
              Yüklediğiniz dosyaya göre bu işlem biraz zaman alabilir.
              Sayfadan ayrılmayın veya yenilemeyin.
            </p>
          </section>
        </div>
      )}
    </>
  );
}