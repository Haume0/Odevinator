import { useEffect, useState } from "react";
import "./App.css";
import { VerifyModal } from "./Components";

function App() {
  const [files, setFiles] = useState([]);
  const [number,setNumber] = useState();
  const [ders,setDers] = useState();
  const [name,setName] = useState()
  const [verifyCode,setVerifyCode] = useState()
  const [verify,setVerify] = useState(false)
  useEffect(() => {
    console.log(files);
  }, [files]);
  // number regex example 2314716027
  const nmbrgx = /^[0-9]{10}$/

function handleSumbit(e) {
    e.preventDefault();
    if(!nmbrgx.test(number)){
      alert("Girdiğiniz okul numarası geçersiz.")
      return
    }
    if(ders == ""){
      alert("Ders adını giriniz.")
      return
    }
    if(files.length == 0){
      alert('Ödevinizi yükleyiniz.')
      return
    }
    if(name == ""){
      alert("Lütfen adınızı giriniz.")
      return
    }
    setVerify(true)
  }
  async function handleFinish(){
    const formData = new FormData();
    for (let i = 0; i < files.length; i++) {
      formData.append("odev_files", files[i]);
    }
    formData.append("ogr_id", number);
    formData.append("ogr_name", name);
    formData.append("ders_name", ders);
    formData.append("verify_code",verifyCode)
    const res = await fetch('http://localhost:8080/odev',{
      method: 'POST',
      body: formData
    })
    console.log(res);
  }
  return (
    <>
      <div className="max-w-[100vw] flex flex-col gap-2 items-center justify-center p-4 md:max-w-[32rem]">
        <img src="/logo.png" className="h-32 aspect-square" alt="" />
        <h1 className="relative bottom-2 font-bold text-5xl">
          Ödevinatör!
        </h1>
        <h4 className="font-medium text-xl">
          Ödevini aşağıdan yükleyebilirsin.
        </h4>
        <form action="" onSubmit={handleSumbit} className="flex flex-col w-full gap-4">
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
              className={`bg-zinc-200 h-12 rounded-xl px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200 ${(!nmbrgx.test(number) && number>0) ? "!border-red-500" : ""}`}
              placeholder="Öğrenci numarası"
              value={number}
              onInput={(e) => setNumber(e.target.value)}
            />
            <input
              type="text"
              className="bg-zinc-200 h-12 rounded-xl px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200"
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
              value={[]}
              className="size-full opacity-0"
            />
            {files.length <= 0 ? (
              <span className="font-semibold absolute pointer-events-none m-auto text-lg">
                Dosyaları eklemek için <br /> buraya bırak / tıkla.
              </span>
            ) : (
              <ul className="flex absolute max-h-full overflow-y-auto py-3 flex-col gap-2">
                {files.map((file, index) => (
                  <div className="flex gap-1">
                    <button
                      type="button"
                      onClick={() => setFiles(files.filter((_, i) => i !== index))}
                      className="px-2 py-1 hover:bg-red-500 text-sm hover:text-white rounded-md">
                      Kaldır
                    </button>
                    <span className="text-lg">{file.name}</span>
                  </div>
                ))}
              </ul>
            )}
          </div>
          <button className="bg-blue-500 active:scale-[0.98] ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
            Gönder
          </button>
        </form>
        <a href="https://bento.me/haume" target="_blank" className="font-medium mt-2 text-zinc-500 text-sm">by Emin Erçoban <span className="text-[10px]">aka</span>Haume</a>
      </div>
      <span>{verifyCode}</span>
      <VerifyModal close={()=>{setVerify(false)}} show={verify} onVerify={(code)=>{setVerifyCode(code)}} />
    </>
  );
}

export default App;
