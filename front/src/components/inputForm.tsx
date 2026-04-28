import RequestPage from "./requestPage";
import { useState } from 'react'
import { useNavigate } from "react-router-dom";
//<button onClick={() => transApply(formData)}>移行申請</button>

type FormData = {
  userID: string
  pageID: string
}

export default function InputForm() {
    const [formData, setFormData] = useState<FormData>({
        userID: "",
        pageID: "",
    })

    const navigate = useNavigate()

    const addData = (e: React.ChangeEvent<HTMLInputElement>) => {
    const key = e.target.name as keyof FormData
    const value = e.target.value

    setFormData((current) => ({
        ...current,
        [key]: value,
    }))
    }

    const goRequestPage = () => {
        if(!formData.userID.trim() || !formData.pageID.trim()){
            alert("userIDとpageIDを入力してください")
            return
        }

        navigate("/request",{
            state: {
                userID: formData.userID,
                pageID: formData.pageID
            }
        })
    }

    return (
    <>
        <h1>Pukiwiki-Migration-Tool</h1>
        <h2>pukiwikiの個人のページのseminar-personal/userIDのuserIDの部分を入力</h2>
        <p>（例：fujimoto2025）</p>
        <input
        name="userID"
        type="text" 
        value={formData.userID}
        onChange={addData}
        placeholder='pukiwikiのuserIDを入力' 
        />

        <h2>ページを追加するnotionのページIDを入力</h2>
        <input
        name="pageID"
        type="text" 
        value={formData.pageID}
        onChange={addData}
        placeholder='notionのpageIDを入力' 
        />
        <br></br>

        
        <p>{formData.userID}</p>
        <p>{formData.pageID}</p>

        <button onClick={goRequestPage}>確認画面へ</button>
        
          
    </>
    )
}