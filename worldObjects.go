package helpers

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	verticies    []float32 //in XYZ UV
	vertexStride int       // 5 if using XYZ UV
	normals      []float32
	bufferLoader *BufferLoader
	vao          BufferID
	nao          BufferID
}

func (o *Object) fillBuffers() {
	o.bufferLoader = NewBufferLoader()
	o.vao = GenBindVertexArray()
	o.nao = GenBindBuffer(gl.ARRAY_BUFFER)

	GenBindBuffer(gl.ARRAY_BUFFER) //VBO

	BindVertexArray(o.vao)
	o.bufferLoader.BuildFloatBuffer(o.vao, NewBufferLayout([]int32{3, 2}, o.verticies))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(o.nao))
	o.bufferLoader.BuildFloatBuffer(o.nao, NewBufferLayout([]int32{3}, o.normals))
}

func (o *Object) calcNormals(triangleCount int) {
	vertexCount := triangleCount * 3 //3 bc we are working in 3d space so XYZ

	o.normals = make([]float32, vertexCount*3)
	for tri := 0; tri < triangleCount; tri++ {
		index := tri * o.vertexStride * 3
		p1 := mgl32.Vec3{o.verticies[index], o.verticies[index+1], o.verticies[index+2]}
		index += o.vertexStride
		p2 := mgl32.Vec3{o.verticies[index], o.verticies[index+1], o.verticies[index+2]}
		index += o.vertexStride
		p3 := mgl32.Vec3{o.verticies[index], o.verticies[index+1], o.verticies[index+2]}

		normal := TriangleNormal(p1, p2, p3)
		o.normals[tri*9+0] = normal.X()
		o.normals[tri*9+1] = normal.Y()
		o.normals[tri*9+2] = normal.Z()

		o.normals[tri*9+3] = normal.X()
		o.normals[tri*9+4] = normal.Y()
		o.normals[tri*9+5] = normal.Z()

		o.normals[tri*9+6] = normal.X()
		o.normals[tri*9+7] = normal.Y()
		o.normals[tri*9+8] = normal.Z()
	}
}

func (o Object) Draw(shader *Shader, drawMatrix mgl32.Mat4) {
	BindVertexArray(o.vao)

	shader.SetMatrix4("model", drawMatrix)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.verticies)/o.vertexStride))
}

func (o Object) DrawMultiple(shader *Shader, num int, drawMatrix func(int) mgl32.Mat4) {
	BindVertexArray(o.vao)

	for i := 0; i < num; i++ {
		shader.SetMatrix4("model", drawMatrix(i))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.verticies)/o.vertexStride))
	}
}

func Cube(size float32) Object {
	o := Object{}
	o.verticies = []float32{
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,
		-size / 2, size / 2, -size / 2, 0.0, 1.0,

		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, size / 2, size / 2, 1.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 1.0,
		-size / 2, size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,

		-size / 2, size / 2, size / 2, 1.0, 0.0,
		-size / 2, size / 2, -size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		-size / 2, size / 2, size / 2, 1.0, 0.0,

		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 0.0,

		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,

		-size / 2, size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		-size / 2, size / 2, -size / 2, 0.0, 1.0,
		-size / 2, size / 2, size / 2, 0.0, 0.0,
	}
	o.vertexStride = 5

	o.calcNormals(12)
	o.fillBuffers()

	return o
}

func Pentahedron(size float32) Object {
	o := Object{}
	o.verticies = []float32{
		size / 2, -size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 0.0,
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		-size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
	}
	o.vertexStride = 5

	o.calcNormals(6)
	o.fillBuffers()

	return o
}
