package ui

func (p *Player) request() {
	select {
	case <-p.perSecond:
		t := int64(p.uiSeekerWidget.Value)
		if p.video == nil {
			select {
			case p.timeChan <- t:
			default:
			}
			break
		}
		ok := t == p.video.Time()
		if !ok {
			select {
			case p.timeChan <- t:
			default:
			}
			break
		}
		canMove := t < p.player.EndTime()
		if !canMove {
			break
		}
		t += 1
		p.uiSeekerWidget.Value = float32(t)
		select {
		case p.timeChan <- t:
		default:
		}

	default:
	}
}
