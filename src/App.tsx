import { A, useLocation, useParams } from "@solidjs/router";
import Icon from "./Icon";
import { useOdevler, useProgress, useUser } from "./Store";
import { createEffect, createSignal, For } from "solid-js";
import { LoadingModal } from "./Loading";

export function AppLayout(props: any) {
  const [user, setUser] = useUser()
  createEffect(() => {
    if (!user.id || user.id == '') {
      window.location.href = '/login'
    }
  })
  return (
    <>
      <div class="bg-gray-300 relative w-[96vw] h-max gap-4 md:w-[480px] p-8 rounded-3xl flex flex-col items-center">
        <header class="w-full flex flex-col gap-4 md:gap-2 md:justify-between">
          <div class="flex items-center gap-2">
            <img src="/logo.png" class=" size-16" alt="" />
            <span class="flex flex-col">
              <h1 class="font-bold text-2xl">Ödevinatör!</h1>
              <p class="text-sm">Ödevleşmek istiyorum.</p>
            </span>
          </div>
          <div class="bg-white gap-2 w-full items-center shrink-0 h-15 flex p-3 rounded-xl">
            <img src={`https://api.dicebear.com/8.x/bottts/svg?seed=${user.id}`} class="size-14 rounded-full" alt="" />
            <span class="flex flex-col">
              <p class="text-sm font-light">Hoşgeldin,</p>
              <h1 class="font-bold text-lg">{user.name} 💠 {user.id}</h1>
            </span>
            <button onClick={() => { setUser({ id: "", code: "", name: "" }); window.open('/login', '_self') }} class="bg-red-500/40 ml-auto size-10 flex group items-center hover:bg-red-500 ease-in-out duration-200 justify-center shrink-0 rounded-lg p-2 text-white">
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
            class={`flex items-center justify-center bg-white rounded-lg shrink-0 hover:bg-blue-500 hover:text-white ease-in-out duration-200 px-4 py-2 font-medium ${useLocation().pathname == "/" ? "!bg-blue-500 text-white" : ""
              }`}>
            Ödev Yükle
          </A>
          <A
            href="/duzenle"
            class={`flex items-center justify-center bg-white rounded-lg shrink-0 hover:bg-blue-500 hover:text-white ease-in-out duration-200 px-4 py-2 font-medium ${useLocation().pathname == "/duzenle" ? "!bg-blue-500 text-white" : ""
              }`}>
            Düzenle
          </A>
        </nav>
        {props.children}
        <a
          href="https://bento.me/haume"
          target="_blank"
          class=" text-blue-500 absolute top-full text-nowrap mx-auto mt-2 hover:underline">
          by Emin “Haume” Erçoban
        </a>
      </div>
      <LoadingModal />
    </>
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
  const [_, setProgress] = useProgress()
  const [user] = useUser()
  const [input, setInput] = createSignal<IInput>({
    lesson: "",
    name: "",
    files: [],
  });
  function handleSumbit(e: any) {
    e.preventDefault();
    setProgress({ state: true, value: 0});
    if (input().lesson == "" || input().name == "" || input().files.length <= 0) {
      alert("Lütfen tüm alanları doldurun.");
      return;
    }
    let formData = new FormData();
    input().files.forEach((file) => {
      formData.append("files", file);
    });
    formData.append("homework_lesson", input().lesson);
    formData.append("homework_name", input().name);
    formData.append("ogr_id", user.id)
    formData.append("ogr_name", user.name)
    formData.append("ogr_code", user.code)
    const xhr = new XMLHttpRequest();

    xhr.upload.onprogress = function (event) {
      setProgress({ state: true, value: (event.loaded / event.total) * 100 });
    };

    xhr.onload = function () {
      if (xhr.status == 200) {
        if(xhr.responseText != "verified"){
          alert("Kullanıcı doğrulanamadı. Lütfen tekrar giriş yapın.");
          setProgress({ state: false, value: 0 });
        }else{
          alert("Ödev başarıyla yüklendi.");
        }
      } else {
        alert("Bir hata oluştu.");
      }
      setProgress({ state: false, value: 0 });
    };

    xhr.onerror = function () {
      alert("Bir hata oluştu.");
      setProgress({ state: false, value: 0 });
    };

    xhr.open("POST", `/new`, true);

    xhr.send(formData);
  }
  return (
    <>
      <h1 class="font-light text-xl md:text-2xl text-center">
        Ödevini aşağıdan yükleyebilirsin.
      </h1>
      <span class="font-bold text-sm md:text-sm text-center">
        Ödev Bilgileri
      </span>
      <form
        action=""
        onSubmit={handleSumbit}
        class="flex flex-col w-full gap-4">
        <input
          type="text"
          class={`bg-white h-12 rounded-xl w-full text-center px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200`}
          placeholder="Dersin adı."
          value={input().lesson}
          onInput={(e) => setInput({ ...input(), lesson: e.target.value })}
        />
        <input
          type="text"
          class={`bg-white h-12 rounded-xl w-full text-center px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200`}
          placeholder="Ödevin adı."
          value={input().name}
          onInput={(e) => setInput({ ...input(), name: e.target.value })}
        />
        <div class="w-full aspect-[2/1] ease-in-out duration-200 rounded-2xl border-dashed relative border-2 hover:border-blue-500 bg-white flex items-center justify-center">
          <input
            onChange={(e) =>
              setInput({
                ...input(),
                files: [...input().files, ...(e.target.files || [])],
              })
            }
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
                      setInput({
                        ...input(),
                        files: [
                          ...input().files.slice(0, index),
                          ...input().files.slice(index + 1),
                        ],
                      })
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
export interface IOdev {
  name: string;
  lesson: string;
  files: IFile[];

}
export function Duzenle() {
  const [user] = useUser()
  const [odevler, setOdevler] = useOdevler()
  createEffect(() => {
    fetch(
      `/odevler?ogr_id=${user.id}&ogr_code=${user.code}`
    )
      .then((res) =>  res.json())
      .then((data) => {
        setOdevler(data);
        // console.log(data)
      });
  })
  return (
    <div class="w-full h-max flex flex-wrap gap-2">
      {odevler.length >0 ? (
        <For each={odevler}>
          {(odev:any,index) => (
            <div class=" relative w-[calc(50%-0.5rem)] flex flex-col gap-1">
            <div class="w-full relative h-full items-center p-2 gap-2 justify-center bg-white rounded-xl flex flex-col">
              <Icon name="folder" class="size-16" />
              <span class="font-medium text-sm md:text-base text-center break-words w-full px-4 leading-4">{odev.name}</span>
              <span class="text-xs">{odev.lesson}</span>
            </div>
            <A href={`/duzenle/${index()}`} class="w-full items-center justify-center flex bg-blue-500 active:scale-[0.96] ease-in-out duration-100 !outline-none text-white py-2 text-sm rounded-lg">
              Düzenle
            </A>
          </div>
          )}
        </For>
      ):(
        <span class="text-center w-full">Henüz yüklenmiş bir ödevin yok.</span>
      )}
    </div>
  );
}
interface IEdit{
  newFiles: File[]
  files: string[]
  lesson: string
  removeFiles: string[]
  name: string
}
export function DuzenleOdev() {
  const [odevler, setOdevler] = useOdevler()
  const index:number =parseInt(useParams().index)
  const [_, setProgress] = useProgress()
  const [user] = useUser()
  // const navigate = useNavigate()
  // console.log(index);
  createEffect(() => {
    
  })
  let oldValues = odevler[index]
  const [input, setInput] = createSignal<IEdit>({newFiles:[],removeFiles:[],files: [], lesson: "", name: ""});
  function handleSumbit(e: any) {
    e.preventDefault();
    setProgress({ state: true, value: 0})
    if (input().lesson == "" || input().name == "") {
      alert("Lütfen alanları boş bırakmayın.");
      return;
    }
    let formData = new FormData();
    input().newFiles.forEach((file) => {
      formData.append("files", file);
    });
    formData.append("homework_lesson", input().lesson);
    formData.append("homework_old_name", oldValues.name);
    formData.append("homework_old_lesson", oldValues.lesson);
    formData.append("homework_name", input().name);
    formData.append("ogr_id", user.id)
    formData.append("ogr_name", user.name)
    formData.append("ogr_code", user.code)
    formData.append("remove_files",JSON.stringify(input().removeFiles))
    const xhr = new XMLHttpRequest();

    xhr.upload.onprogress = function (event) {
      setProgress({ state: true, value: (event.loaded / event.total) * 100 });
    };

    xhr.onload = function () {
      if (xhr.status == 200) {
        if(xhr.responseText != "verified"){
          alert("Kullanıcı doğrulanamadı. Lütfen tekrar giriş yapın.");
          setProgress({ state: false, value: 0 });
        }else{
          alert("Ödev başarıyla düzenlendi.");
        }
      } else {
        alert("Bir hata oluştu.");
      }
      setProgress({ state: false, value: 0 });
      window.location.reload()
      // navigate('/duzenle')
    };
    xhr.onerror = function () {
      alert("Bir hata oluştu.");
      setProgress({ state: false, value: 0 });
      window.location.reload()
      // navigate('/duzenle')
    };

    xhr.open("POST", `/edit`, true);

    xhr.send(formData);
  }
createEffect(() => {
  fetch(
    `/odevler?ogr_id=${user.id}&ogr_code=${user.code}`
  )
    .then((res) => res.json())
    .then((data) => {
      setOdevler(data);
      // console.log(data);
      // @ts-ignore
      oldValues = odevler[index];
      setInput({
        ...input(),
        ...odevler[index],
        //@ts-ignore
        files: odevler[index].files.map((file) => file), // Convert IFile objects to strings
      });
      console.log(input());
      
    });
});
  return (
    <>
      <h1 class="font-light text-xl md:text-2xl text-center">
        Ödevini aşağıdan düzenleyebilirsin.
      </h1>
      <span class="font-bold text-sm md:text-sm text-center">
        Ödev Bilgileri
      </span>
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
          value={input().name}
          onInput={(e) => setInput({ ...input(), name: e.target.value })}
        />
        <div class="w-full aspect-[2/1] ease-in-out duration-200 rounded-2xl border-dashed relative border-2 hover:border-blue-500 bg-white flex items-center justify-center">
          <input
            onChange={(e) =>
              setInput({
                ...input(),
                newFiles: [...input().newFiles, ...(e.target.files || [])],
              })
            }
            type="file"
            accept="*"
            name="files"
            multiple={true}
            class="opacity-0 size-full"
          />
          {input().newFiles.length <= 0 ? (
            <span class="absolute m-auto text-lg text-center font-semibold pointer-events-none">
              Dosyaları eklemek için <br /> buraya bırak / tıkla.
            </span>
          ) : (
            <div class="absolute flex flex-col max-h-full gap-2 py-3 overflow-y-auto">
              {input().newFiles.map((file: any, index) => (
                <div class="flex gap-1">
                  <button
                    type="button"
                    onClick={() =>
                      setInput({
                        ...input(),
                        newFiles: [
                          ...input().newFiles.slice(0, index),
                          ...input().newFiles.slice(index + 1),
                        ],
                      })
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
        <span>Yüklü Dosyalar</span>
          <div class="relative max-h-48 flex flex-col bg-white px-3 rounded-2xl gap-2 py-3 overflow-y-auto">
          {input().files.map((file: string) => (
            <div class="flex gap-1">
              <button
                type="button"
                onClick={() =>      
                  setInput({
                    ...input(),
                    removeFiles: [...input().removeFiles,file],
                    files: input().files.filter((f) => f != file)
                  })
                }
                class="px-2 py-1 text-sm rounded-md shrink-0 bg-zinc-400 hover:bg-red-500 hover:text-white"
                >
                Sil
              </button>
              <div class="overflow-hidden max-w-52 text-nowrap text-ellipsis">
                <span class="text-lg text-ellipsis">{file}</span>
              </div>
            </div>
          ))}
          </div>
        <button class="w-full bg-blue-500 active:scale-[0.98] ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
          Gönder
        </button>
      </form>
    </>
  );
}
