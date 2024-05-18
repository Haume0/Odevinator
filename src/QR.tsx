import { useSearchParams } from "@solidjs/router";
import { createSignal, onMount } from "solid-js";

export default function QR() {
  //get urls from query params
  const [searchParams] = useSearchParams();
  const [urls, setUrls] = createSignal<string[]>([]);
  onMount(() => {
    for (let key in searchParams) {
      if (searchParams.hasOwnProperty(key)) {
        let value = searchParams[key];
        console.log(key, value);
        if (value) {
          setUrls([...urls(), value]);
        }
      }
    }
  });
  return (
    <div class=" flex flex-col items-center justify-center w-[440px] h-full">
      <h1 class="font-bold text-2xl mb-8">Ödevinatör QR kodları!</h1>
      <div class="flex flex-wrap gap-2">
        {urls().map((url) => (
          <div class="flex flex-col items-center gap-2 w-52 p-3 rounded-xl bg-zinc-300">
          <a href={url} target="_blank" class="text-blue-500 max-w-full text-ellipsis overflow-hidden">{url}</a>
          <img src={`https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${url}`} class="w-full rounded-lg aspect-square shrink-0" alt="" />
          </div>
        ))}
      </div>
    </div>
  );
}
