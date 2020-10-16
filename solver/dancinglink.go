package solver

type DancingLink struct {

}

type DancingLinkNode struct {
	cond bool
	left *DancingLinkNode
	right *DancingLinkNode
	up *DancingLinkNode
	down *DancingLinkNode
	
}