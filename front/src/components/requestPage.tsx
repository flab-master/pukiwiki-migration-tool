import { useState } from 'react'
import './requestPage.css'
import { useLocation, useNavigate } from 'react-router-dom'

type Props = {
  userID: string
  pageID: string
}

export default function RequestPage() {
  const location = useLocation()
  const navigate = useNavigate()
  const props = location.state as Props | null
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [message, setMessage] = useState('')

  const transApply = async () => {
    if (!props) {
      navigate('/input')
      return
    }

    if (!props.userID.trim() || !props.pageID.trim()) {
      setMessage('userID と pageID を入力してください。')
      return
    }

    setIsSubmitting(true)
    setMessage('')

    try {
      const response = await fetch('http://localhost:8080/api/migrate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user: props.userID,
          pageID: props.pageID,
        }),
      })

      if (!response.ok) {
        setMessage('移行申請に失敗しました。入力内容を確認してもう一度お試しください。')
        return
      }

      setMessage('移行申請を受け付けました。')
    } catch (error) {
      console.error('通信エラーが発生しました。', error)
      setMessage('通信エラーが発生しました。サーバーの状態を確認してください。')
    } finally {
      setIsSubmitting(false)
    }
  }

  if (!props) {
    return (
      <main className="form-page">
        <section className="panel">
          <div className="page-heading">
            <p className="eyebrow">Step 2</p>
            <h1>入力情報がありません</h1>
            <p>先に移行情報を入力してください。</p>
          </div>
          <button className="primary-button" onClick={() => navigate('/input')}>
            入力画面へ
          </button>
        </section>
      </main>
    )
  }

  return (
    <main className="form-page">
      <section className="panel">
        <div className="page-heading">
          <p className="eyebrow">Step 2</p>
          <h1>申請内容の確認</h1>
          <p>以下の内容で PukiWiki から Notion への移行申請を送信します。</p>
        </div>

        <dl className="confirm-list">
          <div>
            <dt>userID</dt>
            <dd>{props.userID}</dd>
          </div>
          <div>
            <dt>pageID</dt>
            <dd>{props.pageID}</dd>
          </div>
        </dl>

        {message && <p className="status-message">{message}</p>}

        <div className="actions">
          <button className="secondary-button" onClick={() => navigate('/input')}>
            戻る
          </button>
          <button className="primary-button" onClick={transApply} disabled={isSubmitting}>
            {isSubmitting ? '送信中...' : '申請する'}
          </button>
        </div>
      </section>
    </main>
  )
}
