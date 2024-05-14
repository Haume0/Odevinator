import { A, useLocation, useNavigate } from "@solidjs/router";
import Icon from "./Icon";
import { useUser } from "./Store";
import { createEffect, createSignal } from "solid-js";

export function AppLayout(props: { children: Element }) {
  const [user] = useUser();
  return (
    <div class="bg-gray-300 w-[96vw] gap-4 md:w-[480px] p-8 rounded-3xl flex flex-col items-center">
      <header class="w-full flex flex-col gap-4 md:gap-2 md:flex-row md:justify-between">
        <div class="flex items-center gap-2">
          <img src="/logo.png" class=" size-16" alt="" />
          <span class="flex flex-col">
            <h1 class="font-bold text-2xl">Ödevinatör!</h1>
            <p class="text-sm">Ödevleşmek istiyorum.</p>
          </span>
        </div>
        <div class="bg-white gap-2 w-full justify-between md:w-max items-center shrink-0 h-15 flex p-3 rounded-xl">
          <span class="flex flex-col">
            <p class="text-xs font-light">Hoşgeldin,</p>
            <h1 class="font-bold text-base">{user.id}</h1>
          </span>
          <button class="bg-red-500/40 size-10 flex group items-center hover:bg-red-500 ease-in-out duration-200 justify-center shrink-0 rounded-lg p-2 text-white">
            <Icon
              name="power"
              class="size-5 group-hover:text-white ease-in-out duration-200 text-red-500"
            />
          </button>
        </div>
      </header>
      <nav class="flex w-full gap-2 items-center justify-center">
        <A
          href="/"
          class={`flex items-center justify-center bg-white rounded-lg shrink-0 hover:bg-blue-500 hover:text-white ease-in-out duration-200 px-4 py-2 font-medium ${useLocation().pathname == "/" ? "bg-blue-500 text-white" : ""
            }`}>
          Ödev Yükle
        </A>
        <A
          href="/duzenle"
          class={`flex items-center justify-center bg-white rounded-lg shrink-0 hover:bg-blue-500 hover:text-white ease-in-out duration-200 px-4 py-2 font-medium ${useLocation().pathname == "/duzenle" ? "bg-blue-500 text-white" : ""
            }`}>
          Düzenle
        </A>
      </nav>
      {props.children}
      <a href="https://bento.me/haume" target="_blank" class=" ml-auto mt-auto text-blue-500 hover:underline">Created by Emin “Haume” Erçoban</a>
    </div>
  );
}
export interface IFile {
  name: string;
  size: number;
  type: string;
}
export interface IInput {
  lesson: string;
  name: string;
  files: File[];
}
export function App() {
  const [input, setInput] = createSignal<IInput>({
    lesson: "",
    name: "",
    files: [],
  });
  function handleSumbit(e: any) {
    e.preventDefault();
    console.log(input());
  }
  return (
    <>
      <h1 class="font-light text-xl md:text-2xl text-center">Ödevini aşağıdan yükleyebilirsin.</h1>
      <span class="font-bold text-sm md:text-sm text-center">Ödev Bilgileri</span>
      <form
        action=""
        onSubmit={handleSumbit}
        class="flex flex-col w-full gap-4">
        <input
          type="text"
          inputMode="numeric"
          class={`bg-white h-12 rounded-xl w-full text-center px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200`}
          placeholder="Dersin adı."
          value={input().lesson}
          onInput={(e) => setInput({ ...input(), lesson: e.target.value })}
        />
        <input
          type="text"
          inputMode="numeric"
          class={`bg-white h-12 rounded-xl w-full text-center px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200`}
          placeholder="Ödevin adı."
          value={input().lesson}
          onInput={(e) => setInput({ ...input(), lesson: e.target.value })}
        />
        <div class="w-full aspect-[2/1] ease-in-out duration-200 rounded-2xl border-dashed relative border-2 hover:border-blue-500 bg-white flex items-center justify-center">
          <input
            onChange={(e) => setInput({ ...input(), files: [...input().files, ...(e.target.files || [])] })}
            type="file"
            accept="*"
            name="files"
            multiple={true}
            class="opacity-0 size-full"
          />
          {input().files.length <= 0 ? (
            <span class="absolute m-auto text-lg text-center font-semibold pointer-events-none">
              Dosyaları eklemek için <br /> buraya bırak / tıkla.
            </span>
          ) : (
            <div class="absolute flex flex-col max-h-full gap-2 py-3 overflow-y-auto">
              {input().files.map((file: any, index) => (
                <div class="flex gap-1">
                  <button
                    type="button"
                    onClick={() =>
                      setInput({ ...input(), files: [...input().files.slice(0, index), ...input().files.slice(index + 1)] })
                    }
                    class="px-2 py-1 text-sm rounded-md shrink-0 bg-zinc-400 hover:bg-red-500 hover:text-white">
                    Kaldır
                  </button>
                  <div class="overflow-hidden max-w-52 text-nowrap text-ellipsis">
                    <span class="text-lg text-ellipsis">{file.name}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
        <button class="w-full bg-blue-500 active:scale-[0.98] ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
          Gönder
        </button>
      </form>
    </>
  );
}

const odevler = [
  {
    name: "The Pal App",
  },
  {
    name: "The Burgeristan",
  },
  {
    name: "Mackbear Banner Design",
  }
]
export function Duzenle() {
  return (
    <div class="w-full h-max flex flex-wrap gap-2">
      {odevler.map((odev) => (
        <div class=" aspect-square w-[calc(50%-0.5rem)] flex flex-col gap-1">
          <div class="size-full items-center gap-3 justify-center bg-white rounded-xl flex flex-col">
            <Icon name="folder" class="size-16" />
            <span class="font-medium text-center">{odev.name}</span>
          </div>
          <button class="w-full bg-blue-500 active:scale-[0.96] ease-in-out duration-100 !outline-none text-white py-2 text-sm rounded-lg">
          Düzenle
        </button>
        </div>
      ))}
    </div>
  );
}
