import { useEffect, useState } from "react";
import "./App.css";
import { VerifyModal, LoadingModal } from "./Components";
async function postData(url = "", data = {}) {
  // Default options are marked with *
  const response = await fetch(url, {
    method: "POST", // *GET, POST, PUT, DELETE, etc.
    mode: "cors", // no-cors, *cors, same-origin
    cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
    credentials: "same-origin", // include, *same-origin, omit
    headers: {
      "Content-Type": "application/json",
      // 'Content-Type': 'application/x-www-form-urlencoded',
    },
    redirect: "follow", // manual, *follow, error
    referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: JSON.stringify(data), // body data type must match "Content-Type" header
  });
  return response.json(); // parses JSON response into native JavaScript objects
}
function App() {
  const [files, setFiles] = useState([]);
  const [number, setNumber] = useState();
  const [ders, setDers] = useState();
  const [name, setName] = useState();
  const [verify, setVerify] = useState(false);
  const [loading, setLoading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [url,setUrl] = useState(window.location.href);
  if(localStorage.getItem("CLIENT_ID") == null){
    localStorage.setItem("CLIENT_ID", Math.random().toString(36).substring(2, 15))
  }
  // number regex example 2314716027
  const nmbrgx = /^[0-9]{10}$/;
  function handleSumbit(e) {
    e.preventDefault();
    if (!nmbrgx.test(number)) {
      alert("Girdiğiniz okul numarası geçersiz.");
      return;
    }
    if (ders == "") {
      alert("Ders adını giriniz.");
      return;
    }
    if (name == "") {
      alert("Lütfen adınızı giriniz.");
      return;
    }
    postData(`${url}verify`, { ogr_id: number, ogr_name: name, client_id: (localStorage.getItem("CLIENT_ID") || "HATA") })
    setVerify(true);
    if(localStorage.getItem("POSTS") == null){
      localStorage.setItem("POSTS", JSON.stringify([{ogr_id:number,ogr_name:name}])) 
    }else{
      const posts = JSON.parse(localStorage.getItem("POSTS"))
      posts.push({ogr_id:number,ogr_name:name})
      localStorage.setItem("POSTS", JSON.stringify(posts))
    }
  }
  async function handleFinish(code) {
  const formData = new FormData();
  for (let i = 0; i < files.length; i++) {
    formData.append("odev_files", files[i]);
  }
  formData.append("ogr_id", number);
  formData.append("ogr_name", name);
  formData.append("ders_name", ders);
  formData.append("verify_code", code);
  setLoading(true);

  const xhr = new XMLHttpRequest();

  xhr.upload.onprogress = function (event) {
    setProgress((event.loaded / event.total) * 100);
  };

  xhr.onload = function() {
    if (xhr.status == 200) {
      alert("Ödev başarıyla yüklendi.");
      setFiles([]);
      setNumber("");
      setDers("");
      setName("");
    } else {
      alert("Bir hata oluştu.");
    }
    setLoading(false);
  };

  xhr.onerror = function() {
    alert("Bir hata oluştu.");
    setLoading(false);
  };

  xhr.open("POST", `${url}odev`, true);

  xhr.send(formData);
}
  return (
    <>
      <div className="max-w-[100vw] flex flex-col gap-2 items-center justify-center p-4 md:max-w-[32rem]">
        <img src="/logo.png" className="h-32 aspect-square" alt="" />
        <h1 className="relative text-5xl font-bold bottom-2">
          Ödevinatör!
        </h1>
        <h4 className="text-xl font-medium">
          Ödevini aşağıdan yükleyebilirsin.
        </h4>
        <form
          action=""
          onSubmit={handleSumbit}
          className="flex flex-col w-full gap-4">
          <div className="flex flex-col gap-2 mt-4">
            <input
              type="text"
              className={`bg-zinc-200 h-12 rounded-xl px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200`}
              placeholder="Ders adı"
              value={ders}
              onInput={(e) => setDers(e.target.value)}
            />
            <input
              type="text"
              inputMode="numeric"
              className={`bg-zinc-200 h-12 rounded-xl px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200 ${
                !nmbrgx.test(number) && number > 0 ? "!border-red-500" : ""
              }`}
              placeholder="Öğrenci numarası"
              value={number}
              onInput={(e) => setNumber(e.target.value)}
            />
            <input
              type="text"
              className="h-12 px-2 py-2 duration-200 ease-in-out border-2 border-transparent outline-none bg-zinc-200 rounded-xl focus:border-blue-500"
              placeholder="Öğrenci ismi"
              value={name}
              onInput={(e) => setName(e.target.value)}
            />
          </div>
          <div className="h-[14rem] ease-in-out duration-200 aspect-video rounded-2xl border-dashed relative border-2 hover:border-blue-500 bg-zinc-200 flex items-center justify-center">
            <input
              onChange={(e) => setFiles([...e.target.files])}
              type="file"
              multiple="multiple"
              className="opacity-0 size-full"
            />
            {files.length <= 0 ? (
              <span className="absolute m-auto text-lg font-semibold pointer-events-none">
                Dosyaları eklemek için <br /> buraya bırak / tıkla.
              </span>
            ) : (
              <ul className="absolute flex flex-col max-h-full gap-2 py-3 overflow-y-auto">
                {files.map((file, index) => (
                  <div key={index} className="flex gap-1">
                    <button
                      type="button"
                      onClick={() =>
                        setFiles(files.filter((_, i) => i !== index))
                      }
                      className="px-2 py-1 text-sm rounded-md shrink-0 bg-zinc-400 hover:bg-red-500 hover:text-white">
                      Kaldır
                    </button>
                    <div className="overflow-hidden max-w-52 text-nowrap text-ellipsis">
                      <span className="text-lg text-ellipsis">{file.name}</span>
                    </div>
                  </div>
                ))}
              </ul>
            )}
          </div>
          <button className="bg-blue-500 active:scale-[0.98] ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
            Gönder
          </button>
        </form>
        <a
          href="https://bento.me/haume"
          target="_blank"
          className="mt-2 text-sm font-medium text-zinc-500">
          by Emin Erçoban <span className="text-[10px]">aka</span>Haume
        </a>
      </div>
      <VerifyModal
        ogr_id={number}
        close={() => {
          setVerify(false);
        }}
        show={verify}
        onVerify={(code) => {
          handleFinish(code);
        }}
      />
      <LoadingModal show={loading} progress={progress} />
    </>
  );
}

export default App;
