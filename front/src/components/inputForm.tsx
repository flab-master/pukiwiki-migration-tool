import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

type FormData = {
  userID: string
  pageID: string
}

export default function InputForm() {
  const [formData, setFormData] = useState<FormData>({
    userID: '',
    pageID: '',
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

  const goRequestPage = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    if (!formData.userID.trim() || !formData.pageID.trim()) {
      alert('userID と pageID を入力してください。')
      return
    }

    navigate('/request', {
      state: {
        userID: formData.userID.trim(),
        pageID: formData.pageID.trim(),
      },
    })
  }

  return (
    <main className="form-page">
      <section className="panel">
        <div className="page-heading">
          <p className="eyebrow">Step 1</p>
          <h1>移行情報の入力</h1>
          <p>PukiWiki のユーザーIDと、移行先の Notion ページIDを入力してください。</p>
        </div>

        <form className="input-form" onSubmit={goRequestPage}>
          <label className="field">
            <span>PukiWiki userID</span>
            <input
              name="userID"
              type="text"
              value={formData.userID}
              onChange={addData}
              placeholder="例: fujimoto2025"
              autoComplete="username"
            />
            <small>seminar-personal/userID の userID 部分を入力します。</small>
          </label>

          <label className="field">
            <span>Notion pageID</span>
            <input
              name="pageID"
              type="text"
              value={formData.pageID}
              onChange={addData}
              placeholder="Notion のページID"
            />
            <small>移行先として追加する Notion ページのIDです。</small>
          </label>

          <div className="summary-box" aria-label="入力内容">
            <div>
              <span>userID</span>
              <strong>{formData.userID || '未入力'}</strong>
            </div>
            <div>
              <span>pageID</span>
              <strong>{formData.pageID || '未入力'}</strong>
            </div>
          </div>

          <button className="primary-button" type="submit">
            確認画面へ
          </button>
        </form>
      </section>
    </main>
  )
}
