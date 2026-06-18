import { useState, useCallback, useEffect } from 'react'
import './App.css'

const API_BASE = 'http://localhost:8080'

interface Review {
  id: string
  author: string
  score: number
  title: string
  content: string
  submittedAt: string
  appVersion: string
}

type Status = 'idle' | 'loading' | 'error' | 'done'

export default function App() {
  const [appID, setAppID] = useState('595068606')
  const [reviews, setReviews] = useState<Review[]>([])
  const [status, setStatus] = useState<Status>('idle')
  const [error, setError] = useState('')

  const loadReviews = useCallback(async (id: string) => {
    setStatus('loading')
    setError('')
    try {
      const res = await fetch(`${API_BASE}/reviews/${id}`)
      if (!res.ok) throw new Error(`server returned ${res.status}`)
      const data: Review[] | null = await res.json()
      setReviews(data ?? [])
      setStatus('done')
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
      setStatus('error')
    }
  }, [])

  useEffect(() => {
    loadReviews(appID)
  }, [])

  return (
    <div className="page">
      <h1>App Store Reviews</h1>
      <p className="subtitle">Last 720 hours, newest first</p>

      <div className="controls">
        <input
          className="app-input"
          value={appID}
          onChange={(e) => setAppID(e.target.value)}
          placeholder="App ID"
        />
        <button className="load-btn" onClick={() => loadReviews(appID)}>
          Load
        </button>
      </div>

      {status === 'loading' && <p>Loading…</p>}
      {status === 'error' && <p className="feedback error">Error: {error}</p>}
      {status === 'done' && reviews.length === 0 && (
        <p className="feedback empty">No reviews in the last 48 hours.</p>
      )}

      <ul className="review-list">
        {reviews.map((r) => (
          <ReviewCard key={r.id} review={r} />
        ))}
      </ul>
    </div>
  )
}

function ReviewCard({ review }: { review: Review }) {
  return (
    <li className="card">
      <div className="card-header">
        <span className="stars">{renderStars(review.score)}</span>
        <span className="author">{review.author}</span>
        <time className="date" dateTime={review.submittedAt}>
          {new Date(review.submittedAt).toLocaleString()}
        </time>
      </div>
      {review.title && <h3 className="card-title">{review.title}</h3>}
      <p className="card-content">{review.content}</p>
    </li>
  )
}

function renderStars(score: number) {
  return '★'.repeat(score) + '☆'.repeat(5 - score)
}
