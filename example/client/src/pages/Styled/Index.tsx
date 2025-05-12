import styled from 'styled-components'

const StyledExamplePage = () => {
  const Img = styled.img`
  width: 200px;
  height: auto;
  transition: transform 0.3s ease;

  &:hover {
    cursor: pointer;
    transform: scale(1.1); /* Optional: adds a zoom effect on hover */
  }
  `

  return (
    <>
      <Img src="/react.svg" />
      <Img src="/react.svg" />
    </>
  )
}

export { StyledExamplePage }
