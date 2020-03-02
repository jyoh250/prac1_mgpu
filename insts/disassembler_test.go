package insts_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/mgpusim/insts"
)

func TestDisassembler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GCN3 Disassembler")
}

var _ = Describe("Disassembler", func() {

	var (
		disassembler *insts.Disassembler
	)

	BeforeEach(func() {
		disassembler = insts.NewDisassembler()
	})

	It("should decode BF8C0F70", func() {
		buf := []byte{0x70, 0x0f, 0x8c, 0xbf}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("s_waitcnt vmcnt(0)"))
	})

	It("should decode BF8C0171", func() {
		buf := []byte{0x71, 0x01, 0x8c, 0xbf}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("s_waitcnt vmcnt(1) lgkmcnt(1)"))
	})

	It("should decode D81A0004 00000210", func() {
		buf := []byte{0x04, 0x00, 0x1A, 0xd8, 0x10, 0x02, 0x00, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("ds_write_b32 v16, v2 offset:4"))
	})

	It("should decode D86C0008 01000010", func() {
		buf := []byte{0x08, 0x00, 0x6c, 0xd8, 0x10, 0x00, 0x00, 0x01}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("ds_read_b32 v1, v16 offset:8"))
	})

	It("should decode D2850001 00000503", func() {
		buf := []byte{0x01, 0x00, 0x85, 0xd2, 0x03, 0x05, 0x00, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("v_mul_lo_u32 v1, v3, s2"))
	})

	It("should decode D2860004 00000503", func() {
		buf := []byte{0x04, 0x00, 0x86, 0xd2, 0x03, 0x05, 0x00, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("v_mul_hi_u32 v4, v3, s2"))
	})

	It("should decode 041C0D0E", func() {
		buf := []byte{0x0e, 0x0d, 0x1c, 0x04}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("v_sub_f32_e32 v14, v14, v6"))
	})

	It("should decode 7E224911", func() {
		buf := []byte{0x11, 0x49, 0x22, 0x7e}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("v_rsq_f32_e32 v17, v17"))
	})

	It("should decode 7E540900", func() {
		buf := []byte{0x00, 0x09, 0x54, 0x7e}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("v_cvt_f64_i32_e32 v[42:43], v0"))
	})

	It("should decode 7E221F0F", func() {
		buf := []byte{0x0f, 0x1f, 0x22, 0x7e}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).To(Equal("v_cvt_f32_f64_e32 v17, v[15:16]"))
	})

	It("should decode D281000F 0000012A", func() {
		buf := []byte{0x0f, 0x00, 0x81, 0xd2, 0x2a, 0x01, 0x00, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_mul_f64 v[15:16], v[42:43], s[0:1]"))
	})

	It("should decode D04E0100 00000111", func() {
		buf := []byte{0x00, 0x01, 0x4e, 0xd0, 0x11, 0x01, 0x00, 0x00}

		inst, err := disassembler.Decode(buf)

		Expect(err).To(BeNil())
		Expect(inst.String(nil)).
			To(Equal("v_cmp_nlt_f32_e64 s[0:1], |v17|, s0"))
	})
})
