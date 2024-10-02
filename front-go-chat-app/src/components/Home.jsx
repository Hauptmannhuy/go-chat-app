import { useEffect, useState } from "react"



function Home(){
  const [data, setData] = useState(null)
  const [isLoading, setLoading] = useState(true)
  
  

    useEffect(() => {
      const fetchData = async () => {
      // setLoading(true)
      const response = await fetch('http://localhost:8090/')
      const message = await response.json()
      setData(message)
      setLoading(false)
    }
    fetchData()
  }, [])
  console.log(data)

  
  if (isLoading){
    return (<><h2>...loading</h2></>)
  }

  return (<>
    <h2>{data.Message}</h2>
  
  </>
  )
}

export default Home