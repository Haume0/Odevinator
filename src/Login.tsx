import { useNavigate } from "@solidjs/router";
import { createSignal, Match, Switch } from "solid-js";
import { useUser } from "./Store";

export function LoginLayout(props: any) {
  return (
    <div class="bg-gray-300 w-[96vw] h-max gap-4 md:w-[480px] p-8 rounded-3xl flex flex-col items-center">
      <header class="w-full flex flex-col gap-4 md:gap-2 md:flex-row items-center md:justify-between">
        <div class="flex items-center gap-2">
          <img src="/logo.png" class=" size-16" alt="" />
          <span class="flex flex-col">
            <h1 class="font-bold text-2xl">Ödevinatör!</h1>
            <p class="text-sm">Ödevleşmek istiyorum.</p>
          </span>
        </div>
      </header>
      {props.children}
      <a
        href="https://bento.me/haume"
        target="_blank"
        class=" ml-auto mt-auto text-blue-500 hover:underline">
        by Emin “Haume” Erçoban
      </a>
    </div>
  );
}

export function Login() {
  // 2314716027 -> 10 character long id regex
  const idrgx = /^[0-9]{10}$/;
  const [_,setUser] = useUser();
  const [input, setInput] = createSignal({
    name: "",
    id: "",
    code: "",
  });
  const navigate = useNavigate();
  const [stage, setStage] = createSignal<"id" | "code">("id");
  function handleLogin(e: any) {
    e.preventDefault();
    if (stage() == "id") {
      if (!idrgx.test(input().id)) {
        alert("Lütfen geçerli bir okul numarası girin.");
        return;
      }
      //fetch /login?id=${input().id}
      fetch(
        `/auth?id=${input().id}&name=${input().name}`
      ).then(async (res) => {
        if ((await res.text()) == "exists") {
          alert("Zaten kayıtlısınız!\nYeni kod gönderilmeyecek!\nLütfen e-postanı kontrol et.");
        } else if ((await res.text()) == "done") {
          alert("E-postanıza doğrulama kodu gönderildi!");
        }
      });
      setStage("code");
    } else if (stage() == "code") {
      //fetch /verify?code=${input().code}
      fetch(
        `/verify?id=${input().id}&code=${input().code}`
      ).then(async (res) => {
        if ((await res.text()) == "verified") {
          alert("Giriş başarılı! Yönlendiriliyorsunuz..");
          setUser({ id: input().id, code: input().code, name: input().name});
          navigate("/");
        }else{
          alert("Kayıt bulunamadı ya da kod geçersiz!");
        }
      });
    }
  }
  return (
    <Switch>
      <Match when={stage() == "id"}>
        <>
          <h1 class=" text-4xl text-center">Giriş Yap!</h1>
          <p class="font-medium text-center">Lütfen okul numaranızı girin.</p>
          <form
            onSubmit={handleLogin}
            action=""
            class="flex flex-col w-full gap-4">
            <input
              type="text"
              class={`bg-white h-12 rounded-xl w-full text-center px-2 py-2 outline-none focus:!border-blue-500 border-2 border-transparent ease-in-out duration-200`}
              placeholder="İsminiz"
              value={input().name}
              onInput={(e) => setInput({ ...input(), name: e.target.value })}
            />
            <input
              type="text"
              class={`bg-white h-12 rounded-xl w-full text-center px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200 ${
                idrgx.test(input().id) ? "border-green-500" : "!border-red-500"
              }`}
              placeholder="Okul numaranız"
              inputMode="numeric"
              value={input().id}
              onInput={(e) => setInput({ ...input(), id: e.target.value })}
            />
            <button class="w-full bg-blue-500 active:scale-[0.98] ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
              Gönder
            </button>
          </form>
        </>
      </Match>
      <Match when={stage() == "code"}>
        <>
          <h1 class=" text-4xl text-center">Doğrula</h1>
          <p class="font-medium text-center">
            Lütfen <mark>okul e-postanıza</mark> gelen <br />6 haneli kodu
            girin.
          </p>
          <form
            onSubmit={handleLogin}
            action=""
            class="flex flex-col w-full gap-4">
            <input
              type="text"
              class={`bg-white h-12 tracking-[0.5em] rounded-xl w-full text-center px-2 py-2 outline-none focus:border-blue-500 border-2 border-transparent ease-in-out duration-200`}
              placeholder="XXXXXX"
              value={input().code}
              onInput={(e) => setInput({ ...input(), code: e.target.value })}
            />
            <div class="flex w-full gap-2">
              <button
                onClick={() => {
                  setStage("id");
                }}
                class="w-1/3 bg-zinc-500 active:scale-[0.98] ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
                Geri
              </button>
              <button class="w-full bg-blue-500 active:scale-[0.98] ease-in-out duration-100 !outline-none text-white font-bold py-2 px-4 rounded-lg">
                Gönder
              </button>
            </div>
          </form>
        </>
      </Match>
    </Switch>
  );
}
