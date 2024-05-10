import React from "react";
import { useState } from "react";

export function VerifyModal(props: { show: boolean; ogr_id:string; close: () => void; onVerify: (code: string) => void }) {
  const [code, setCode] = useState("");
  function handleVerify() {
    fetch(`http://localhost:8080/check?code=${code}&id=${props.ogr_id}`).then((res)=>res.json()).then((data)=>{
      if(data.msg != "OK"){
        alert("Doğrulama kodu yanlış!")
        return
      }
      alert('Kod doğru!')
      props.onVerify(code);
      props.close();
    }).catch((err)=>{
      console.log(err)
      alert('Bir sorunla karşılaştık!')
    })
  }
  return (
    <>
      {props.show && (
        <div className=" fixed bg-black/60 backdrop-blur-sm w-screen flex items-center justify-center h-[100svh] inset-0">
          <section className="w-[32rem] p-4 bg-white rounded-xl flex flex-col gap-4">
            <h1 className="text-3xl font-bold">Lütfen sen olduğunu doğrula!</h1>
            <p className="font-medium">
              Sen olduğunu doğrulamak için OKUL E-POSTA adresine gelen 6 haneli
              kodu gir.
            </p>
            <input
              type="text"
              className={`bg-zinc-200 h-12 text-center rounded-xl px-4 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200`}
              placeholder="Doğrulama kodu"
              value={code}
              onInput={(e) => setCode(e.target.value)}
            />
            <div className="flex w-full gap-4">
              <button
                onClick={() => props.close()}
                className="bg-zinc-500 active:scale-[0.98] w-1/3 ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
                Geri
              </button>
              <button
                onClick={handleVerify}
                className="bg-blue-500 active:scale-[0.98] w-full ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
                Doğrula
              </button>
            </div>
          </section>
        </div>
      )}
    </>
  );
}

export function LoadingModal(props: { show: boolean }) {
  return(
    <>
      {props.show && (
        <div className=" fixed bg-black/60 backdrop-blur-sm w-screen flex items-center justify-center h-[100svh] inset-0">
          <section className="w-[32rem] p-12 bg-white items-center justify-center rounded-xl flex flex-col gap-4">
            <img src="/spinner.gif" className="size-24" alt="" />
            <h1 className="text-3xl font-bold">Lütfen bekleyin...</h1>
            <p className="font-medium">
              Yüklediğiniz dosyaya göre bu işlem biraz zaman alabilir.
              Sayfadan ayrılmayın veya yenilemeyin.
            </p>
          </section>
        </div>
      )}
    </>
  );
}