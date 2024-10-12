import { useEffect, useState } from "react"
import { Link } from "react-router-dom";


function Home(){
  const [data, setData] = useState(null)
  // const [isLoading, setLoading] = useState(true)
  
  

    useEffect(() => {
      const fetchData = async () => {
      const response = await fetch('http://localhost:8090/sign_up')
      const message = await response.json()
      // setData(message)
      // setLoading(false)
    }
    fetchData()
  }, [])
  console.log(data)

  
  // if (isLoading){
  //   return (<><h2>...loading</h2></>)
  // }

  return (<>
   
  </>
  )
}

export default Home