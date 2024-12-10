import { useRef } from "react"

export default function SearchSection({onSearch}){
  const timeout = useRef(null)
 
  function typeEventHandler(input){
    if (timeout.current){
      clearTimeout(timeout.current)
    }
    timeout.current = setTimeout(()=>{
      onSearch(input)
    }, 400)
  }
  return (<>
  <input type="text" className="searchInput" onChange={(e) => {typeEventHandler(e.target.value)}} />
  </>)
}