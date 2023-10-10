package commcid_test

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"testing"

	commcid "github.com/filecoin-project/go-fil-commcid"
	commhash "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multihash"
	"github.com/multiformats/go-varint"
	"github.com/stretchr/testify/require"
)

func TestDataCommitmentToCID(t *testing.T) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)

	c, err := commcid.DataCommitmentV1ToCID(randBytes)
	require.NoError(t, err)

	require.Equal(t, c.Prefix().Codec, uint64(cid.FilCommitmentUnsealed))
	mh := c.Hash()
	decoded, err := multihash.Decode([]byte(mh))
	require.NoError(t, err)
	require.Equal(t, decoded.Code, uint64(multihash.SHA2_256_TRUNC254_PADDED))
	require.Equal(t, decoded.Length, len(randBytes))
	require.True(t, bytes.Equal(decoded.Digest, randBytes))
}

func TestCIDToDataCommitment(t *testing.T) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)

	t.Run("with correct hash format", func(t *testing.T) {
		hash := testMultiHash(multihash.SHA2_256_TRUNC254_PADDED, randBytes, 0)

		t.Run("decodes raw commitment hash when correct cid format", func(t *testing.T) {
			c := cid.NewCidV1(cid.FilCommitmentUnsealed, hash)
			decoded, err := commcid.CIDToDataCommitmentV1(c)
			require.NoError(t, err)
			require.True(t, bytes.Equal(decoded, randBytes))
		})

		t.Run("error on non-fil codec", func(t *testing.T) {
			c := cid.NewCidV1(cid.DagCBOR, hash)
			decoded, err := commcid.CIDToDataCommitmentV1(c)
			require.EqualError(t, err, commcid.ErrIncorrectCodec.Error())
			require.Nil(t, decoded)
		})

		t.Run("error on wrong fil codec", func(t *testing.T) {
			c := cid.NewCidV1(cid.FilCommitmentSealed, testMultiHash(multihash.POSEIDON_BLS12_381_A1_FC1, randBytes, 0))
			decoded, err := commcid.CIDToDataCommitmentV1(c)
			require.EqualError(t, err, commcid.ErrIncorrectCodec.Error())
			require.Nil(t, decoded)
		})

		t.Run("error on fil hash/codec mismatch", func(t *testing.T) {
			c := cid.NewCidV1(cid.FilCommitmentUnsealed, testMultiHash(multihash.POSEIDON_BLS12_381_A1_FC1, randBytes, 0))
			decoded, err := commcid.CIDToDataCommitmentV1(c)
			require.EqualError(t, err, commcid.ErrIncorrectHash.Error())
			require.Nil(t, decoded)
		})

	})

	t.Run("error on incorrectly formatted hash", func(t *testing.T) {
		hash := testMultiHash(multihash.SHA2_256_TRUNC254_PADDED, randBytes, 5)
		c := cid.NewCidV1(cid.FilCommitmentUnsealed, hash)
		decoded, err := commcid.CIDToDataCommitmentV1(c)
		require.Error(t, err)
		require.Regexp(t, "^Error decoding data commitment hash:", err.Error())
		require.Nil(t, decoded)
	})
}

func TestReplicaCommitmentToCID(t *testing.T) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)

	c, err := commcid.ReplicaCommitmentV1ToCID(randBytes)
	require.NoError(t, err)

	require.Equal(t, c.Prefix().Codec, uint64(cid.FilCommitmentSealed))
	mh := c.Hash()
	decoded, err := multihash.Decode([]byte(mh))
	require.NoError(t, err)
	require.Equal(t, decoded.Code, uint64(multihash.POSEIDON_BLS12_381_A1_FC1))
	require.Equal(t, decoded.Length, len(randBytes))
	require.True(t, bytes.Equal(decoded.Digest, randBytes))
}

func TestCIDToReplicaCommitment(t *testing.T) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)

	t.Run("with correct hash format", func(t *testing.T) {
		hash := testMultiHash(multihash.POSEIDON_BLS12_381_A1_FC1, randBytes, 0)

		t.Run("decodes raw commitment hash when correct cid format", func(t *testing.T) {
			c := cid.NewCidV1(cid.FilCommitmentSealed, hash)
			decoded, err := commcid.CIDToReplicaCommitmentV1(c)
			require.NoError(t, err)
			require.True(t, bytes.Equal(decoded, randBytes))
		})

		t.Run("error on incorrect CID format", func(t *testing.T) {
			c := cid.NewCidV1(cid.DagCBOR, hash)
			decoded, err := commcid.CIDToReplicaCommitmentV1(c)
			require.EqualError(t, err, commcid.ErrIncorrectCodec.Error())
			require.Nil(t, decoded)
		})

		t.Run("error on non-fil codec", func(t *testing.T) {
			c := cid.NewCidV1(cid.DagCBOR, hash)
			decoded, err := commcid.CIDToReplicaCommitmentV1(c)
			require.EqualError(t, err, commcid.ErrIncorrectCodec.Error())
			require.Nil(t, decoded)
		})

		t.Run("error on wrong fil codec", func(t *testing.T) {
			c := cid.NewCidV1(cid.FilCommitmentUnsealed, testMultiHash(multihash.SHA2_256_TRUNC254_PADDED, randBytes, 0))
			decoded, err := commcid.CIDToReplicaCommitmentV1(c)
			require.EqualError(t, err, commcid.ErrIncorrectCodec.Error())
			require.Nil(t, decoded)
		})

		t.Run("error on fil hash/codec mismatch", func(t *testing.T) {
			c := cid.NewCidV1(cid.FilCommitmentSealed, testMultiHash(multihash.SHA2_256_TRUNC254_PADDED, randBytes, 0))
			decoded, err := commcid.CIDToReplicaCommitmentV1(c)
			require.EqualError(t, err, commcid.ErrIncorrectHash.Error())
			require.Nil(t, decoded)
		})
	})

	t.Run("error on wrong hash type", func(t *testing.T) {
		encoded, err := multihash.Encode(randBytes, multihash.SHA2_256)
		require.NoError(t, err)
		c := cid.NewCidV1(cid.Raw, multihash.Multihash(encoded))
		decoded, err := commcid.CIDToReplicaCommitmentV1(c)
		require.EqualError(t, err, commcid.ErrIncorrectCodec.Error())
		require.Nil(t, decoded)
	})

	t.Run("error on incorrectly formatted hash", func(t *testing.T) {
		c := cid.NewCidV1(cid.FilCommitmentUnsealed, testMultiHash(multihash.POSEIDON_BLS12_381_A1_FC1, randBytes, 5))
		decoded, err := commcid.CIDToReplicaCommitmentV1(c)
		require.Error(t, err)
		require.Regexp(t, "^Error decoding data commitment hash:", err.Error())
		require.Nil(t, decoded)
	})

}

func TestPieceCommitmentToCID(t *testing.T) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)

	c, err := commcid.PieceCommitmentV1ToCID(randBytes)
	require.NoError(t, err)

	require.Equal(t, c.Prefix().Codec, uint64(cid.FilCommitmentUnsealed))
	mh := c.Hash()
	decoded, err := multihash.Decode([]byte(mh))
	require.NoError(t, err)
	require.Equal(t, decoded.Code, uint64(multihash.SHA2_256_TRUNC254_PADDED))
	require.Equal(t, decoded.Length, len(randBytes))
	require.True(t, bytes.Equal(decoded.Digest, randBytes))

	_, err = commcid.PieceCommitmentV1ToCID(randBytes[1:])
	require.Regexp(t, "^commitments must be 32 bytes long", err.Error())
}

func TestCIDToPieceCommitment(t *testing.T) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)

	t.Run("with correct hash format", func(t *testing.T) {
		hash := testMultiHash(multihash.SHA2_256_TRUNC254_PADDED, randBytes, 0)

		t.Run("decodes raw commitment hash when correct cid format", func(t *testing.T) {
			c := cid.NewCidV1(cid.FilCommitmentUnsealed, hash)
			decoded, err := commcid.CIDToPieceCommitmentV1(c)
			require.NoError(t, err)
			require.True(t, bytes.Equal(decoded, randBytes))
		})

		t.Run("error on incorrect CID format", func(t *testing.T) {
			c := cid.NewCidV1(cid.DagCBOR, hash)
			decoded, err := commcid.CIDToPieceCommitmentV1(c)
			require.EqualError(t, err, commcid.ErrIncorrectCodec.Error())
			require.Nil(t, decoded)
		})
	})

	t.Run("error on incorrectly formatted hash", func(t *testing.T) {
		hash := testMultiHash(multihash.SHA2_256_TRUNC254_PADDED, randBytes, 5)
		c := cid.NewCidV1(cid.FilCommitmentUnsealed, hash)
		decoded, err := commcid.CIDToPieceCommitmentV1(c)
		require.Error(t, err)
		require.Regexp(t, "^Error decoding data commitment hash:", err.Error())
		require.Nil(t, decoded)
	})
	t.Run("error on wrong hash type", func(t *testing.T) {
		encoded, err := multihash.Encode(randBytes, multihash.SHA2_256)
		require.NoError(t, err)
		c := cid.NewCidV1(cid.FilCommitmentUnsealed, multihash.Multihash(encoded))
		decoded, err := commcid.CIDToPieceCommitmentV1(c)
		require.EqualError(t, err, commcid.ErrIncorrectHash.Error())
		require.Nil(t, decoded)
	})
}

func randomPieceMhInfo(t *testing.T) (treeHeight uint8, paddingSize uint64, dataSize uint64, digest []byte, mhDigest []byte) {
	// CID size = 1 byte tree height + unsigned_varint_data_padding_size (1 to 9 bytes) + 32 byte digest
	// Min/Max size = 34 -> 42
	t.Helper()

	randBytes := make([]byte, 40)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)
	digest = randBytes[0:32]
	dataSize = binary.LittleEndian.Uint64(randBytes[32:])
	// padded dataSize must be less than 2^63 - 1 so we divide by 4 to be safe
	dataSize = dataSize >> 2
	// TODO: at the moment some of the code here requires max size to be 128 times lower
	dataSize = dataSize >> 7
	// minimum dataSize is 127
	dataSize += 127
	treeHeight, paddingSize, err = commcid.UnpaddedSizeToV1TreeHeightAndPadding(dataSize)
	require.NoError(t, err)

	uvarintPaddingSize := varint.ToUvarint(paddingSize)
	mhDigest = append(append(uvarintPaddingSize[:], treeHeight), digest...)
	return
}

func TestPieceCommitmentToPieceMhCID(t *testing.T) {
	height, paddingSize, dataSize, digest, _ := randomPieceMhInfo(t)

	c, err := commcid.DataCommitmentV1ToPieceMhCID(digest, dataSize)
	require.NoError(t, err)

	require.Equal(t, c.Prefix().Codec, uint64(cid.Raw))
	mh := c.Hash()
	decoded, err := multihash.Decode([]byte(mh))
	require.NoError(t, err)
	require.Equal(t, decoded.Code, uint64(commcid.FR32_SHA256_TRUNC254_PADDED_BINARY_TREE_CODE))
	require.Equal(t, decoded.Length, 1+varint.UvarintSize(paddingSize)+32)
	require.True(t, decoded.Digest[varint.UvarintSize(paddingSize)] == height)

	paddingSizeFromMhDigest, _, err := varint.FromUvarint(decoded.Digest[0:varint.UvarintSize(paddingSize)])
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, paddingSize, paddingSizeFromMhDigest)

	_, err = commcid.DataCommitmentV1ToPieceMhCID(digest[1:], dataSize)
	require.Regexp(t, "^commitments must be 32 bytes long", err.Error())
}

func TestPieceMhCIDToPieceCommitment(t *testing.T) {
	treeHeight, paddingSize, expectedDataSize, expectedDigest, mhDigest := randomPieceMhInfo(t)

	t.Run("with correct hash format", func(t *testing.T) {
		hash := testMultiHash(commcid.FR32_SHA256_TRUNC254_PADDED_BINARY_TREE_CODE, mhDigest, 0)

		t.Run("decodes raw commitment hash when correct cid format", func(t *testing.T) {
			c := cid.NewCidV1(cid.Raw, hash)
			digest, dataSize, err := commcid.PieceMhCIDToDataCommitmentV1(c)
			t.Log(treeHeight)
			t.Log(paddingSize)
			require.NoError(t, err)
			require.Equal(t, expectedDataSize, dataSize)
			require.True(t, bytes.Equal(expectedDigest, digest))
		})

		t.Run("don't error on non-Raw CID format", func(t *testing.T) {
			c := cid.NewCidV1(cid.DagCBOR, hash)
			digest, dataSize, err := commcid.PieceMhCIDToDataCommitmentV1(c)
			require.NoError(t, err)
			require.Equal(t, expectedDataSize, dataSize)
			require.True(t, bytes.Equal(expectedDigest, digest))
		})
	})

	t.Run("error on incorrectly formatted hash", func(t *testing.T) {
		hash := testMultiHash(commcid.FR32_SHA256_TRUNC254_PADDED_BINARY_TREE_CODE, mhDigest, 5)
		c := cid.NewCidV1(cid.Raw, hash)
		digest, _, err := commcid.PieceMhCIDToDataCommitmentV1(c)
		require.Error(t, err)
		require.Regexp(t, "^Error decoding data commitment hash:", err.Error())
		require.Nil(t, digest)
	})
	t.Run("error on wrong hash type", func(t *testing.T) {
		encoded, err := multihash.Encode(mhDigest, multihash.SHA2_256)
		require.NoError(t, err)
		c := cid.NewCidV1(cid.Raw, multihash.Multihash(encoded))
		digest, _, err := commcid.PieceMhCIDToDataCommitmentV1(c)
		require.EqualError(t, err, commcid.ErrIncorrectHash.Error())
		require.Nil(t, digest)
	})
}

func TestTreeHeight(t *testing.T) {
	// Add test fixtures
	noFr32PaddingTests := map[string]struct {
		size    uint64
		height  uint8
		padding int64
	}{
		"127OfEach0-1-2-3":          {127 * 4, 4, 0},
		"512-bytes-should-pad-over": {512, 5, 504},
		"0":                         {0, 0, -1},
		"1":                         {1, 0, -1},
		"31":                        {31, 0, -1},
		"32":                        {32, 1, -1},
		"127":                       {127, 2, 0},
		"32GiB":                     {32 << 30, 31, 33822867456},
		"32GiB-post-padding":        {(32 << 30) * 127 / 128, 30, 0},
		"64GiB":                     {64 << 30, 32, 67645734912},
		"64GiB-post-padding":        {(64 << 30) * 127 / 128, 31, 0},
	}

	for name, tc := range noFr32PaddingTests {
		t.Run(fmt.Sprintf("non-fr32-padding %s", name), func(t *testing.T) {
			t.Run("height-only", func(t *testing.T) {
				height, err := commcid.UnpaddedSizeToV1TreeHeight(tc.size)
				require.NoError(t, err)
				require.Equal(t, tc.height, height)
			})
			if tc.size >= 127 {
				t.Run("height-and-padding", func(t *testing.T) {
					height, padding, err := commcid.UnpaddedSizeToV1TreeHeightAndPadding(tc.size)
					require.NoError(t, err)
					require.Equal(t, tc.height, height)
					require.Equal(t, uint64(tc.padding), padding)
				})
			}
		})
	}

	// Add test fixtures
	fr32PaddingTests := map[string]struct {
		size   uint64
		height uint8
	}{
		"127OfEach0-1-2-3":              {127 * 4, 4},
		"512-bytes-should-not-pad-over": {512, 4},
		"0":                             {0, 0},
		"1":                             {1, 0},
		"31":                            {31, 0},
		"32":                            {32, 0},
		"127":                           {127, 2},
		"128":                           {128, 2},
		"129":                           {129, 3},
		"32GiB":                         {32 << 30, 30},
		"64GiB":                         {64 << 30, 31},
	}

	for name, tc := range fr32PaddingTests {
		t.Run(fmt.Sprintf("with-fr32-padding %s", name), func(t *testing.T) {
			require.Equal(t, tc.height, commcid.Fr32PaddedSizeToV1TreeHeight(tc.size))
		})
	}
}

func TestMultihashes(t *testing.T) {
	data127EachOf0_1_2_3 := append(append(append(bytes.Repeat([]byte{0x00}, 127), bytes.Repeat([]byte{0x01}, 127)...), bytes.Repeat([]byte{0x02}, 127)...), bytes.Repeat([]byte{0x03}, 127)...)

	// Add test fixtures
	tests := map[string]struct {
		data     []byte
		v2CidStr string
	}{
		"127OfEach0-1-2-3":              {data127EachOf0_1_2_3, "bafkzcibcaaces3nobte6ezpp4wqan2age2s5yxcatzotcvobhgcmv5wi2xh5mbi"},
		"127OfEach0-1-2-3-Then127*4-0s": {append(data127EachOf0_1_2_3[:], bytes.Repeat([]byte{0x00}, 127*4)...), "bafkzcibcaac542av3szurbbscwuu3zjssvfwbpsvbjf6y3tukvlgl2nf5rha6pa"},
		"127OfEach0-1-2-3-Then127+4-0s": {append(data127EachOf0_1_2_3[:], bytes.Repeat([]byte{0x00}, 127+4)...), "bafkzcibd7ebalxticxolgseegik2stpfgkkuwyf6kufex3doorkvmzpjuxwe4dz4"},
		"127OfEach0-1-2-3-Then127+5-0s": {append(data127EachOf0_1_2_3[:], bytes.Repeat([]byte{0x00}, 127+5)...), "bafkzcibd7abalxticxolgseegik2stpfgkkuwyf6kufex3doorkvmzpjuxwe4dz4"},
	}

	for name, tc := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			v2Cid, err := cid.Parse(tc.v2CidStr)
			require.NoError(t, err)

			h := &commhash.Calc{}
			_, err = h.Write(tc.data)
			if err != nil {
				t.Fatal(err)
			}
			digest, paddedSize, err := h.Digest()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("digest: %x", digest)
			t.Logf("unpadded size: %d", len(tc.data))
			t.Logf("padded fr32 data size: %d", paddedSize)

			computedV2Cid, err := commcid.DataCommitmentV1ToPieceMhCID(digest, uint64(len(tc.data)))
			require.NoError(t, err)

			cidStr, err := computedV2Cid.StringOfBase(multibase.Base16)
			if err != nil {
				t.Fatal(err)
			}

			c, _, err := commcid.ConvertDataCommitmentV1PieceMhCIDToV1CID(computedV2Cid)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf(cidStr)
			t.Logf(computedV2Cid.String())
			t.Logf(c.String())

			require.Equal(t, v2Cid, computedV2Cid)
		})
	}

}

func TestPieceMhCIDandV1CIDPieceCommitmentConverters(t *testing.T) {
	_, _, expectedDataSize, expectedDigest, mhDigest := randomPieceMhInfo(t)

	mhv1 := testMultiHash(multihash.SHA2_256_TRUNC254_PADDED, expectedDigest, 0)
	cidv1 := cid.NewCidV1(cid.FilCommitmentUnsealed, mhv1)

	mhv2 := testMultiHash(commcid.FR32_SHA256_TRUNC254_PADDED_BINARY_TREE_CODE, mhDigest, 0)
	cidv2 := cid.NewCidV1(cid.Raw, mhv2)

	t.Run("convert v1 piece cid + data size to piece mh cid", func(t *testing.T) {
		c, err := commcid.ConvertDataCommitmentV1V1CIDtoPieceMhCID(cidv1, expectedDataSize)
		require.NoError(t, err)
		require.Equal(t, cidv2, c)
	})

	t.Run("convert piece mh cid to v1 piece cid + data size", func(t *testing.T) {
		c, dataSize, err := commcid.ConvertDataCommitmentV1PieceMhCIDToV1CID(cidv2)
		require.NoError(t, err)
		require.Equal(t, cidv1, c)
		require.Equal(t, expectedDataSize, dataSize)
	})

	// Add test fixtures
	tests := map[string]struct {
		v1CidStr         string
		unpaddedDataSize uint64
		v2CidStr         string
	}{
		"127OfEach0-1-2-3":              {"baga6ea4seaqes3nobte6ezpp4wqan2age2s5yxcatzotcvobhgcmv5wi2xh5mbi", 127 * 4, "bafkzcibcaaces3nobte6ezpp4wqan2age2s5yxcatzotcvobhgcmv5wi2xh5mbi"},
		"empty32GiB":                    {"baga6ea4seaqao7s73y24kcutaosvacpdjgfe5pw76ooefnyqw4ynr3d2y6x2mpq", (32 << 30) * 127 / 128, "bafkzcibcaapao7s73y24kcutaosvacpdjgfe5pw76ooefnyqw4ynr3d2y6x2mpq"},
		"empty64GiB":                    {"baga6ea4seaqomqafu276g53zko4k23xzh4h4uecjwicbmvhsuqi7o4bhthhm4aq", (64 << 30) * 127 / 128, "bafkzcibcaap6mqafu276g53zko4k23xzh4h4uecjwicbmvhsuqi7o4bhthhm4aq"},
		"127OfEach0-1-2-3-Then127*4-0s": {"baga6ea4seaqn42av3szurbbscwuu3zjssvfwbpsvbjf6y3tukvlgl2nf5rha6pa", 127 * 8, "bafkzcibcaac542av3szurbbscwuu3zjssvfwbpsvbjf6y3tukvlgl2nf5rha6pa"},
		"127OfEach0-1-2-3-Then127+4-0s": {"baga6ea4seaqn42av3szurbbscwuu3zjssvfwbpsvbjf6y3tukvlgl2nf5rha6pa", 127*4 + 127 + 4, "bafkzcibd7ebalxticxolgseegik2stpfgkkuwyf6kufex3doorkvmzpjuxwe4dz4"},
		"127OfEach0-1-2-3-Then127+5-0s": {"baga6ea4seaqn42av3szurbbscwuu3zjssvfwbpsvbjf6y3tukvlgl2nf5rha6pa", 127*4 + 127 + 5, "bafkzcibd7abalxticxolgseegik2stpfgkkuwyf6kufex3doorkvmzpjuxwe4dz4"},
	}

	for name, tc := range tests {
		t.Run(fmt.Sprintf("%s-v1-to-v2", name), func(t *testing.T) {
			v1Cid, err := cid.Parse(tc.v1CidStr)
			require.NoError(t, err)

			v2Cid, err := cid.Parse(tc.v2CidStr)
			require.NoError(t, err)

			computedV2Cid, err := commcid.ConvertDataCommitmentV1V1CIDtoPieceMhCID(v1Cid, tc.unpaddedDataSize)
			require.NoError(t, err)

			require.Equal(t, v2Cid, computedV2Cid)
		})

		t.Run(fmt.Sprintf("%s-v2-to-v1", name), func(t *testing.T) {
			v1Cid, err := cid.Parse(tc.v1CidStr)
			require.NoError(t, err)

			v2Cid, err := cid.Parse(tc.v2CidStr)
			require.NoError(t, err)

			computedV1Cid, computedHeight, err := commcid.ConvertDataCommitmentV1PieceMhCIDToV1CID(v2Cid)
			require.NoError(t, err)

			require.Equal(t, v1Cid, computedV1Cid)
			require.Equal(t, tc.unpaddedDataSize, computedHeight)
		})
	}
}

func testMultiHash(code uint64, buf []byte, extra int) multihash.Multihash {
	newBuf := make([]byte, varint.UvarintSize(code)+varint.UvarintSize(uint64(len(buf)))+len(buf)+extra)
	n := varint.PutUvarint(newBuf, code)
	n += varint.PutUvarint(newBuf[n:], uint64(len(buf)))

	copy(newBuf[n:], buf)
	return multihash.Multihash(newBuf)
}
