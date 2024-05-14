import { createSignal } from "solid-js"

export default function Icon(props: { name: string, class?: string, width?: string, height?: string }) {
  let [icon,setIcon] = createSignal<string>("") // Icon variable
  fetch(`/${props.name}.svg`).then(async(a) => {
    setIcon(await a.text()) // Converting svg file to text
  }).then(() => {
    let dummy:string = icon()
    //change all fills to currentcolor
    dummy = dummy.replace(/fill="[^"]*"/g, 'fill="currentColor"')
    dummy = dummy.replace("<svg", `<svg class="${props.class}" width="${props.width}" height="${props.height}"`)
    setIcon(dummy)
  })
  //CSS
  const css: any = `
  #cubicon { /* Core css rules */
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}
  `
  return (
    <>
      {icon() != '' &&
        <div id="cubicon" style={css} class={props.class} innerHTML={icon()}></div>
      }
    </>
  )
}